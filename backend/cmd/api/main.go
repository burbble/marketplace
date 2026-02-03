package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/burbble/marketplace/internal/config"
	"github.com/burbble/marketplace/pkg/db"
)

const (
	shutdownTimeout = 5 * time.Second
	startTimeout    = 30 * time.Second
	stopTimeout     = 30 * time.Second
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
	env       = config.EnvDev
)

func main() {
	fx.New(
		fx.Provide(
			ProvideConfig,
			ProvideLogger,
			ProvideDB,
			ProvideRedis,
			ProvideRouter,
			ProvideHTTPServer,
		),
		fx.Invoke(SetupRoutes),
		fx.Invoke(StartServer),

		fx.StartTimeout(startTimeout),
		fx.StopTimeout(stopTimeout),
	).Run()
}

func ProvideConfig() (*config.Config, error) {
	cfg := &config.Config{}
	if err := config.LoadFromFlags(cfg); err != nil {
		return nil, err
	}

	cfg.Environment = env

	return cfg, nil
}

func ProvideLogger(cfg *config.Config) (*zap.Logger, error) {
	var lg *zap.Logger
	var err error

	if cfg.IsDevEnv() {
		lg, err = zap.NewDevelopment()
	} else {
		lg, err = zap.NewProduction()
	}

	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	lg = lg.With(
		zap.String("version", version),
		zap.String("env", env),
	)

	_ = zap.ReplaceGlobals(lg)

	return lg, nil
}

func ProvideDB(ctx context.Context, cfg *config.Config, lg *zap.Logger) (*db.Connection, error) {
	conn, err := db.NewConnection(ctx, &cfg.PostgresConfig, lg)
	if err != nil {
		lg.Error("failed to connect to postgres", zap.Error(err))
		return nil, err
	}

	return conn, nil
}

func ProvideRedis(ctx context.Context, cfg *config.Config, lg *zap.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		lg.Error("failed to connect to redis", zap.Error(err))
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	lg.Info("redis connected", zap.String("addr", cfg.Addr()))

	return rdb, nil
}

func ProvideRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

func ProvideHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
	}
}

func SetupRoutes(
	cfg *config.Config,
	router *gin.Engine,
	lg *zap.Logger,
) {
	apiV1 := router.Group("/api/v1")

	apiV1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": version})
	})

	lg.Info("routes registered")
}

func StartServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	srv *http.Server,
	lg *zap.Logger,
	conn *db.Connection,
	rdb *redis.Client,
) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				lg.Info("starting server", zap.String("addr", srv.Addr))
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					lg.Error("server error", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			lg.Info("shutdown signal received")

			shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				lg.Error("server forced to shutdown", zap.Error(err))
				return err
			}

			rdb.Close()
			conn.Close()

			lg.Info("server stopped")

			return nil
		},
	})
}
