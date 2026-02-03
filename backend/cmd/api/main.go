package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/burbble/marketplace/docs"
	"github.com/burbble/marketplace/internal/config"
	"github.com/burbble/marketplace/internal/exchange"
	"github.com/burbble/marketplace/internal/handler"
	"github.com/burbble/marketplace/internal/repository/postgres"
	"github.com/burbble/marketplace/internal/service"
	"github.com/burbble/marketplace/pkg/db"
	"github.com/burbble/marketplace/pkg/ratelimit"
	"github.com/burbble/marketplace/pkg/zapx"
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

// @title          Store Marketplace API
// @version        1.0
// @description    Product catalog API for store77.net marketplace
// @BasePath       /api/v1
func main() {
	fx.New(
		fx.Provide(
			ProvideConfig,
			ProvideLogger,
			ProvideDB,
			ProvideRedis,
			ProvideRouter,
			ProvideHTTPServer,
			postgres.NewCategoryRepo,
			postgres.NewProductRepo,
			service.NewCategoryService,
			service.NewProductService,
			exchange.NewGrinexProvider,
			handler.NewCategoryHandler,
			handler.NewProductHandler,
			handler.NewExchangeHandler,
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
	return zapx.Init(cfg.LogMode,
		zap.String("service", "api"),
		zap.String("version", version),
		zap.String("env", env),
	)
}

func ProvideDB(cfg *config.Config, lg *zap.Logger) (*db.Connection, error) {
	ctx := context.Background()

	conn, err := db.NewConnection(ctx, &cfg.PostgresConfig, lg)
	if err != nil {
		lg.Error("failed to connect to postgres", zap.Error(err))
		return nil, err
	}

	return conn, nil
}

func ProvideRedis(cfg *config.Config, lg *zap.Logger) (*redis.Client, error) {
	ctx := context.Background()

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

func ProvideRouter(cfg *config.Config, rdb *redis.Client) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware(),
		ratelimit.Middleware(rdb, ratelimit.Config{
			Max:    cfg.RateLimitRPS,
			Window: time.Second,
		}),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Header("Access-Control-Max-Age", "43200")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func SetupRoutes(
	router *gin.Engine,
	lg *zap.Logger,
	ch *handler.CategoryHandler,
	ph *handler.ProductHandler,
	eh *handler.ExchangeHandler,
) {
	apiV1 := router.Group("/api/v1")

	apiV1.GET("/categories", ch.List)
	apiV1.GET("/categories/:id", ch.GetByID)

	apiV1.GET("/products", ph.List)
	apiV1.GET("/products/:id", ph.GetByID)

	apiV1.GET("/brands", ph.GetBrands)

	apiV1.GET("/exchange/rate", eh.GetRate)

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
