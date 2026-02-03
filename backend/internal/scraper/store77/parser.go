package store77

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Category struct {
	Name string
	URL  string
	Slug string
}

type Product struct {
	ExternalID  string
	SKU         string
	Name        string
	Price       int
	ImageURL    string
	ProductURL  string
	Brand       string
	Category    string
	Description string
}

type PaginationInfo struct {
	CurrentPage int
	TotalPages  int
}

type ecomProduct struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Price    any    `json:"price"`
	Brand    string `json:"brand"`
	Category string `json:"category"`
}

var (
	ecomRe    = regexp.MustCompile(`YandexEcommerce\.getInstance\(\)\.\w+\(\[(.+?)\]\)`)
	spacesRe  = regexp.MustCompile(`\s+`)
)

func normalizeSpace(s string) string {
	return strings.TrimSpace(spacesRe.ReplaceAllString(s, " "))
}

func ParseCategories(html string) ([]Category, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	seen := make(map[string]struct{})
	var categories []Category

	doc.Find("ul.catalog_menu > li").Each(func(_ int, topLi *goquery.Selection) {
		topLi.Find("ul.catalog_menu_sub_second > li").Each(func(_ int, secLi *goquery.Selection) {
			thirdUl := secLi.Find("ul.catalog_menu_sub_third")

			if thirdUl.Length() > 0 {
				thirdUl.Find("li > a").Each(func(_ int, a *goquery.Selection) {
					href, _ := a.Attr("href")
					name := strings.TrimSpace(a.Text())
					addCategory(&categories, seen, name, href)
				})
			} else {
				a := secLi.Find("div.bli_pos_second > a")
				href, _ := a.Attr("href")
				name := strings.TrimSpace(a.Text())
				addCategory(&categories, seen, name, href)
			}
		})
	})

	return categories, nil
}

func addCategory(categories *[]Category, seen map[string]struct{}, name, href string) {
	if href == "" || name == "" || href == "#" {
		return
	}

	href = strings.TrimSpace(href)
	name = normalizeSpace(name)

	if _, ok := seen[href]; ok {
		return
	}
	seen[href] = struct{}{}

	slug := strings.Trim(href, "/")
	if idx := strings.Index(slug, "?"); idx != -1 {
		slug = slug[:idx]
	}

	*categories = append(*categories, Category{
		Name: name,
		URL:  href,
		Slug: slug,
	})
}

func ParseProducts(html string) ([]Product, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	var products []Product

	doc.Find("div.wrap_list_prod div.blocks_product").Each(func(_ int, card *goquery.Selection) {
		p := Product{}

		if btn := card.Find("button.favorite_product"); btn.Length() > 0 {
			p.ExternalID, _ = btn.Attr("data-elid")
		}

		link := card.Find("div.blocks_product_fix_w > a").First()
		if link.Length() == 0 {
			return
		}

		href, _ := link.Attr("href")
		p.ProductURL = href

		if img := link.Find("img"); img.Length() > 0 {
			p.ImageURL, _ = img.Attr("src")
			if p.Name == "" {
				p.Name = strings.TrimSpace(img.AttrOr("title", ""))
			}
		}

		if h2 := card.Find("h2.bp_text_info a"); h2.Length() > 0 {
			p.Name = normalizeSpace(h2.Text())
		}

		if priceEl := card.Find("p.bp_text_price"); priceEl.Length() > 0 {
			p.Price = parsePrice(priceEl.Text())
		}

		onclick, exists := link.Attr("onclick")
		if exists {
			if ecom := parseEcomData(onclick); ecom != nil {
				p.SKU = ecom.ID
				p.Brand = ecom.Brand
				p.Category = ecom.Category
				if p.Price == 0 {
					p.Price = ecomPrice(ecom.Price)
				}
			}
		}

		if p.Name != "" && p.ProductURL != "" {
			products = append(products, p)
		}
	})

	return products, nil
}

func ParsePagination(html string) (*PaginationInfo, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	info := &PaginationInfo{
		CurrentPage: 1,
		TotalPages:  1,
	}

	pag := doc.Find("div.pagination_catalog ul.pagination").First()
	if pag.Length() == 0 {
		return info, nil
	}

	if active := pag.Find("a.active"); active.Length() > 0 {
		if n, err := strconv.Atoi(strings.TrimSpace(active.Text())); err == nil {
			info.CurrentPage = n
		}
	}

	pag.Find("li > a").Each(func(_ int, a *goquery.Selection) {
		text := strings.TrimSpace(a.Text())
		if text == "..." || text == "" {
			return
		}
		if n, err := strconv.Atoi(text); err == nil && n > info.TotalPages {
			info.TotalPages = n
		}
	})

	return info, nil
}

func ParseTotalProducts(html string) (int, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return 0, fmt.Errorf("parse html: %w", err)
	}

	text := doc.Find("div.search_all_produkt span").First().Text()
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, nil
	}

	n, err := strconv.Atoi(strings.ReplaceAll(text, " ", ""))
	if err != nil {
		return 0, fmt.Errorf("parse total: %w", err)
	}

	return n, nil
}

func ParseProductDescription(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	selectors := []string{
		"div.detail_text",
		"div#detail_text",
		"div[itemprop='description']",
		"div.product-detail-text",
		"div.product_description",
		"div.element-detail-text",
	}

	for _, sel := range selectors {
		el := doc.Find(sel).First()
		if el.Length() > 0 {
			text := strings.TrimSpace(el.Text())
			if text != "" {
				return normalizeSpace(text)
			}
		}
	}

	return ""
}

func parsePrice(s string) int {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\u00a0", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.TrimRight(s, "—-–")
	s = strings.TrimSpace(s)

	n, _ := strconv.Atoi(s)
	return n
}

func parseEcomData(onclick string) *ecomProduct {
	matches := ecomRe.FindStringSubmatch(onclick)
	if len(matches) < 2 {
		return nil
	}

	jsonStr := matches[1]

	var p ecomProduct
	if err := json.Unmarshal([]byte(jsonStr), &p); err != nil {
		return nil
	}

	return &p
}

func ecomPrice(v any) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case string:
		n, _ := strconv.Atoi(strings.ReplaceAll(val, " ", ""))
		return n
	}
	return 0
}
