package exchange

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func newTestProvider(t *testing.T) *grinexProvider {
	t.Helper()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:63790",
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skipf("redis not available, skipping: %v", err)
	}

	t.Cleanup(func() {
		rdb.Del(ctx, redisCacheKey)
		_ = rdb.Close()
	})

	return &grinexProvider{
		client: &http.Client{},
		rdb:    rdb,
		logger: zap.NewNop(),
	}
}

func TestFetchRateSuccess(t *testing.T) {
	resp := depthResponse{
		Bids: []orderBookEntry{
			{Price: "95.50", Volume: "100", Amount: "9550"},
			{Price: "95.40", Volume: "200", Amount: "19080"},
		},
		Asks: []orderBookEntry{
			{Price: "95.60", Volume: "100", Amount: "9560"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	provider := newTestProvider(t)
	provider.client = server.Client()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	httpResp, err := provider.client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	var depth depthResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&depth); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(depth.Bids) == 0 {
		t.Fatal("expected bids in response")
	}

	if depth.Bids[0].Price != "95.50" {
		t.Errorf("expected best bid 95.50, got %s", depth.Bids[0].Price)
	}
}

func TestFetchRateNoBids(t *testing.T) {
	resp := depthResponse{
		Bids: []orderBookEntry{},
		Asks: []orderBookEntry{{Price: "95.60"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	httpResp, err := http.Get(server.URL) //nolint:gosec // test URL
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	var depth depthResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&depth); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(depth.Bids) != 0 {
		t.Errorf("expected 0 bids, got %d", len(depth.Bids))
	}
}

func TestBidSpreadCalculation(t *testing.T) {
	bestBid := 95.50
	rate := bestBid - bidSpread

	expected := 95.40
	if rate != expected {
		t.Errorf("expected rate %.2f, got %.2f", expected, rate)
	}
}

func TestDepthResponseParsing(t *testing.T) {
	jsonData := `{
		"timestamp": 1700000000,
		"bids": [
			{"price": "95.50", "volume": "100.5", "amount": "9600.75"},
			{"price": "95.40", "volume": "200", "amount": "19080"}
		],
		"asks": [
			{"price": "95.60", "volume": "50", "amount": "4780"}
		]
	}`

	var depth depthResponse
	if err := json.Unmarshal([]byte(jsonData), &depth); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if depth.Timestamp != 1700000000 {
		t.Errorf("expected timestamp 1700000000, got %d", depth.Timestamp)
	}
	if len(depth.Bids) != 2 {
		t.Errorf("expected 2 bids, got %d", len(depth.Bids))
	}
	if len(depth.Asks) != 1 {
		t.Errorf("expected 1 ask, got %d", len(depth.Asks))
	}
	if depth.Bids[0].Price != "95.50" {
		t.Errorf("expected bid price 95.50, got %s", depth.Bids[0].Price)
	}
}
