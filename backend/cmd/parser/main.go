package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/burbble/marketplace/internal/config"
	"github.com/burbble/marketplace/pkg/db"
	"github.com/burbble/marketplace/pkg/zapx"
)

type exitCode = int

const (
	noErr exitCode = iota
	errLoadConfig
	errInitLogger
	errConnectToDB
	errConnectToRedis
	errScrape
)

type application struct {
	logger *zap.Logger
	cfg    *config.Config
	conn   *db.Connection
	rdb    *redis.Client
}

func main() {
	os.Exit(run())
}

func run() exitCode {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := &config.Config{}
	if err := config.LoadFromFlags(cfg); err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		return errLoadConfig
	}

	lg, err := zapx.Init(cfg.LogMode, zap.String("service", "parser"))
	if err != nil {
		fmt.Printf("failed to init logger: %v\n", err)
		return errInitLogger
	}

	conn, err := db.NewConnection(ctx, &cfg.PostgresConfig, lg)
	if err != nil {
		lg.Error("failed to connect to postgres", zap.Error(err))
		return errConnectToDB
	}
	defer conn.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		lg.Error("failed to connect to redis", zap.Error(err))
		return errConnectToRedis
	}

	lg.Info("redis connected", zap.String("addr", cfg.Addr()))

	app := &application{
		logger: lg,
		cfg:    cfg,
		conn:   conn,
		rdb:    rdb,
	}

	return app.runScraper(ctx)
}

func (a *application) runScraper(ctx context.Context) exitCode {
	a.logger.Info("starting scraper", zap.Duration("interval", a.cfg.ScrapeInterval))

	ticker := time.NewTicker(a.cfg.ScrapeInterval)
	defer ticker.Stop()

	if err := a.scrape(ctx); err != nil {
		a.logger.Error("scrape failed", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("scraper stopped")
			return noErr
		case <-ticker.C:
			if err := a.scrape(ctx); err != nil {
				a.logger.Error("scrape failed", zap.Error(err))
			}
		}
	}
}

func (a *application) scrape(ctx context.Context) error {
	a.logger.Info("scraping started")

	a.logger.Info("scraping completed")

	return nil
}
