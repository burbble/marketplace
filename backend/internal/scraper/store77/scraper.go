package store77

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"go.uber.org/zap"
)

const (
	baseURL         = "https://store77.net"
	pageLoadTimeout = 30 * time.Second
	httpTimeout     = 15 * time.Second
)

type Scraper struct {
	logger     *zap.Logger
	browser    *rod.Browser
	httpClient *http.Client
}

func NewScraper(logger *zap.Logger) *Scraper {
	return &Scraper{
		logger: logger,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func (s *Scraper) Start() error {
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).
		Headless(true).
		Set("disable-gpu").
		Set("no-sandbox").
		MustLaunch()

	s.browser = rod.New().ControlURL(u)
	if err := s.browser.Connect(); err != nil {
		return fmt.Errorf("connect to browser: %w", err)
	}

	s.logger.Info("browser started")

	return nil
}

func (s *Scraper) Stop() {
	if s.browser != nil {
		_ = s.browser.Close()
		s.logger.Info("browser closed")
	}
}

func (s *Scraper) FetchPageHTML(ctx context.Context, url string) (string, error) {
	page := s.browser.MustPage()
	defer page.MustClose()

	page = page.Context(ctx).Timeout(pageLoadTimeout)

	err := page.Navigate(url)
	if err != nil {
		return "", fmt.Errorf("navigate to %s: %w", url, err)
	}

	err = page.WaitLoad()
	if err != nil {
		return "", fmt.Errorf("wait load %s: %w", url, err)
	}

	time.Sleep(3 * time.Second)

	html, err := page.HTML()
	if err != nil {
		return "", fmt.Errorf("get html from %s: %w", url, err)
	}

	return html, nil
}

func (s *Scraper) FetchMainPage(ctx context.Context) (string, error) {
	s.logger.Info("fetching main page", zap.String("url", baseURL))

	html, err := s.FetchPageHTML(ctx, baseURL)
	if err != nil {
		return "", err
	}

	s.logger.Info("main page fetched", zap.Int("html_length", len(html)))

	return html, nil
}

func (s *Scraper) FetchProductPage(ctx context.Context, path string) (string, error) {
	u := baseURL + path

	s.logger.Debug("fetching product page", zap.String("url", u))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("create request %s: %w", u, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch product page %s: %w", u, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("product page %s returned status %d", u, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read product page %s: %w", u, err)
	}

	return string(body), nil
}

func (s *Scraper) FetchCategoryPage(ctx context.Context, path string, page int) (string, error) {
	u := baseURL + path
	if page > 1 {
		u += fmt.Sprintf("?PAGEN_1=%d", page)
	}

	s.logger.Info("fetching category page", zap.String("url", u), zap.Int("page", page))

	html, err := s.FetchPageHTML(ctx, u)
	if err != nil {
		return "", err
	}

	s.logger.Info("category page fetched", zap.Int("html_length", len(html)))

	return html, nil
}
