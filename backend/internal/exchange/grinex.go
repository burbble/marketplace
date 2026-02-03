package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	grinexDepthURL = "https://grinex.io/api/v1/spot/depth?symbol=usdta7a5"
	redisCacheKey  = "exchange:usdt_rub"
	cacheTTL       = 1 * time.Minute
	httpTimeout    = 10 * time.Second
	bidSpread      = 0.10
)

type RateProvider interface {
	GetUSDTRate(ctx context.Context) (float64, error)
}

type orderBookEntry struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
}

type depthResponse struct {
	Timestamp int64            `json:"timestamp"`
	Bids      []orderBookEntry `json:"bids"`
	Asks      []orderBookEntry `json:"asks"`
}

type grinexProvider struct {
	client *http.Client
	rdb    *redis.Client
	logger *zap.Logger
}

func NewGrinexProvider(rdb *redis.Client, logger *zap.Logger) RateProvider {
	return &grinexProvider{
		client: &http.Client{Timeout: httpTimeout},
		rdb:    rdb,
		logger: logger,
	}
}

func (g *grinexProvider) GetUSDTRate(ctx context.Context) (float64, error) {
	cached, err := g.rdb.Get(ctx, redisCacheKey).Float64()
	if err == nil {
		return cached, nil
	}

	rate, err := g.fetchRate(ctx)
	if err != nil {
		return 0, err
	}

	if err := g.rdb.Set(ctx, redisCacheKey, rate, cacheTTL).Err(); err != nil {
		g.logger.Warn("failed to cache exchange rate", zap.Error(err))
	}

	return rate, nil
}

func (g *grinexProvider) fetchRate(ctx context.Context) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, grinexDepthURL, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetch grinex depth: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("grinex returned status %d", resp.StatusCode)
	}

	var depth depthResponse
	if err := json.NewDecoder(resp.Body).Decode(&depth); err != nil {
		return 0, fmt.Errorf("decode grinex response: %w", err)
	}

	if len(depth.Bids) == 0 {
		return 0, fmt.Errorf("no bids in grinex response")
	}

	bestBid, err := strconv.ParseFloat(depth.Bids[0].Price, 64)
	if err != nil {
		return 0, fmt.Errorf("parse bid price %q: %w", depth.Bids[0].Price, err)
	}

	rate := bestBid - bidSpread

	g.logger.Info("fetched exchange rate",
		zap.Float64("best_bid", bestBid),
		zap.Float64("rate", rate),
	)

	return rate, nil
}
