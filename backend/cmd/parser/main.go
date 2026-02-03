package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/burbble/marketplace/internal/config"
	"github.com/burbble/marketplace/internal/domain"
	"github.com/burbble/marketplace/internal/repository/postgres"
	"github.com/burbble/marketplace/internal/scraper/store77"
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
	logger       *zap.Logger
	cfg          *config.Config
	conn         *db.Connection
	rdb          *redis.Client
	scraper      *store77.Scraper
	categoryRepo postgres.CategoryRepository
	productRepo  postgres.ProductRepository
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
		logger:       lg,
		cfg:          cfg,
		conn:         conn,
		rdb:          rdb,
		scraper:      store77.NewScraper(lg),
		categoryRepo: postgres.NewCategoryRepo(conn),
		productRepo:  postgres.NewProductRepo(conn),
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

	if err := a.scraper.Start(); err != nil {
		return fmt.Errorf("start browser: %w", err)
	}
	defer a.scraper.Stop()

	mainHTML, err := a.scraper.FetchMainPage(ctx)
	if err != nil {
		return fmt.Errorf("fetch main page: %w", err)
	}

	parsedCategories, err := store77.ParseCategories(mainHTML)
	if err != nil {
		return fmt.Errorf("parse categories: %w", err)
	}

	a.logger.Info("categories parsed", zap.Int("count", len(parsedCategories)))

	seen := make(map[string]struct{}, len(parsedCategories))
	domainCategories := make([]domain.Category, 0, len(parsedCategories))
	for _, c := range parsedCategories {
		if _, exists := seen[c.Slug]; exists {
			continue
		}
		seen[c.Slug] = struct{}{}
		domainCategories = append(domainCategories, domain.Category{
			Name: c.Name,
			Slug: c.Slug,
			URL:  c.URL,
		})
	}

	if err := a.categoryRepo.Upsert(ctx, domainCategories); err != nil {
		return fmt.Errorf("upsert categories: %w", err)
	}

	dbCategories, err := a.categoryRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("get categories: %w", err)
	}

	slugToID := make(map[string]uuid.UUID, len(dbCategories))
	for _, c := range dbCategories {
		slugToID[c.Slug] = c.ID
	}

	workers := a.cfg.ScrapeWorkers
	if workers <= 0 {
		workers = 5
	}

	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	a.logger.Info("scraping categories", zap.Int("concurrency", workers), zap.Int("total", len(parsedCategories)))

	for _, cat := range parsedCategories {
		categoryID, ok := slugToID[cat.Slug]
		if !ok {
			continue
		}

		select {
		case <-ctx.Done():
			break
		case sem <- struct{}{}:
		}

		wg.Add(1)
		go func(cat store77.Category, categoryID uuid.UUID) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := a.scrapeCategory(ctx, cat, categoryID); err != nil {
				a.logger.Error("scrape category failed",
					zap.String("category", cat.Name),
					zap.Error(err),
				)
			}
		}(cat, categoryID)
	}

	wg.Wait()

	a.logger.Info("scraping completed")
	return nil
}

func (a *application) scrapeCategory(ctx context.Context, cat store77.Category, categoryID uuid.UUID) error {
	a.logger.Info("scraping category", zap.String("name", cat.Name), zap.String("url", cat.URL))

	html, err := a.scraper.FetchCategoryPage(ctx, cat.URL, 1)
	if err != nil {
		return fmt.Errorf("fetch page 1: %w", err)
	}

	pagination, err := store77.ParsePagination(html)
	if err != nil {
		return fmt.Errorf("parse pagination: %w", err)
	}

	a.logger.Info("category pagination",
		zap.String("category", cat.Name),
		zap.Int("total_pages", pagination.TotalPages),
	)

	if err := a.processPage(ctx, html, categoryID); err != nil {
		return fmt.Errorf("process page 1: %w", err)
	}

	for page := 2; page <= pagination.TotalPages; page++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		pageHTML, err := a.scraper.FetchCategoryPage(ctx, cat.URL, page)
		if err != nil {
			a.logger.Error("fetch page failed",
				zap.String("category", cat.Name),
				zap.Int("page", page),
				zap.Error(err),
			)
			continue
		}

		if err := a.processPage(ctx, pageHTML, categoryID); err != nil {
			a.logger.Error("process page failed",
				zap.String("category", cat.Name),
				zap.Int("page", page),
				zap.Error(err),
			)
			continue
		}
	}

	return nil
}

func (a *application) processPage(ctx context.Context, html string, categoryID uuid.UUID) error {
	parsed, err := store77.ParseProducts(html)
	if err != nil {
		return fmt.Errorf("parse products: %w", err)
	}

	if len(parsed) == 0 {
		return nil
	}

	products := make([]domain.Product, 0, len(parsed))
	for _, p := range parsed {
		if p.ExternalID == "" {
			continue
		}

		price := p.Price - 1000
		if price < 0 {
			price = 0
		}

		products = append(products, domain.Product{
			ExternalID:    p.ExternalID,
			SKU:           p.SKU,
			Name:          p.Name,
			OriginalPrice: p.Price,
			Price:         price,
			ImageURL:      p.ImageURL,
			ProductURL:    p.ProductURL,
			Brand:         p.Brand,
			CategoryID:    categoryID,
		})
	}

	if len(products) == 0 {
		return nil
	}

	a.logger.Info("upserting products", zap.Int("count", len(products)))

	return a.productRepo.Upsert(ctx, products)
}
