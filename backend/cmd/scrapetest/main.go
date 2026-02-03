package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/burbble/marketplace/internal/scraper/store77"
)

func main() {
	lg, _ := zap.NewDevelopment()

	scraper := store77.NewScraper(lg)
	if err := scraper.Start(); err != nil {
		lg.Fatal("failed to start scraper", zap.Error(err))
	}
	defer scraper.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// 1. fetch main page and parse categories
	html, err := scraper.FetchMainPage(ctx)
	if err != nil {
		lg.Fatal("failed to fetch main page", zap.Error(err))
	}
	_ = os.WriteFile("store77_dump.html", []byte(html), 0644)

	categories, err := store77.ParseCategories(html)
	if err != nil {
		lg.Fatal("failed to parse categories", zap.Error(err))
	}

	fmt.Printf("\n=== CATEGORIES (%d) ===\n", len(categories))
	for i, c := range categories {
		if i >= 20 {
			fmt.Printf("  ... and %d more\n", len(categories)-20)
			break
		}
		fmt.Printf("  [%d] %s → %s\n", i+1, c.Name, c.URL)
	}

	catURL := "https://store77.net/telefony_apple/"
	catHTML, err := scraper.FetchPageHTML(ctx, catURL)
	if err != nil {
		lg.Fatal("failed to fetch category page", zap.Error(err))
	}
	_ = os.WriteFile("store77_category_dump.html", []byte(catHTML), 0644)

	products, err := store77.ParseProducts(catHTML)
	if err != nil {
		lg.Fatal("failed to parse products", zap.Error(err))
	}

	total, _ := store77.ParseTotalProducts(catHTML)
	pag, _ := store77.ParsePagination(catHTML)

	fmt.Printf("\n=== PRODUCTS (page %d/%d, total %d) ===\n", pag.CurrentPage, pag.TotalPages, total)
	for i, p := range products {
		if i >= 10 {
			fmt.Printf("  ... and %d more on this page\n", len(products)-10)
			break
		}
		fmt.Printf("  [%d] %s | %d руб | SKU: %s | Brand: %s\n", i+1, p.Name, p.Price, p.SKU, p.Brand)
		fmt.Printf("       URL: %s\n", p.ProductURL)
		fmt.Printf("       Category: %s\n", p.Category)
	}
}
