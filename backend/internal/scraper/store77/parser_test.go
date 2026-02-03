package store77

import (
	"testing"
)

func TestParseCategories(t *testing.T) {
	html := `<html><body>
		<ul class="catalog_menu">
			<li>
				<ul class="catalog_menu_sub_second">
					<li>
						<div class="bli_pos_second"><a href="/phones/">Телефоны</a></div>
					</li>
					<li>
						<ul class="catalog_menu_sub_third">
							<li><a href="/phones/apple/">Apple</a></li>
							<li><a href="/phones/samsung/">Samsung</a></li>
						</ul>
					</li>
				</ul>
			</li>
		</ul>
	</body></html>`

	cats, err := ParseCategories(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cats) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(cats))
	}

	if cats[0].Name != "Телефоны" {
		t.Errorf("expected name 'Телефоны', got %q", cats[0].Name)
	}
	if cats[0].URL != "/phones/" {
		t.Errorf("expected URL '/phones/', got %q", cats[0].URL)
	}
	if cats[0].Slug != "phones" {
		t.Errorf("expected slug 'phones', got %q", cats[0].Slug)
	}

	if cats[1].Name != "Apple" {
		t.Errorf("expected name 'Apple', got %q", cats[1].Name)
	}
	if cats[1].Slug != "phones/apple" {
		t.Errorf("expected slug 'phones/apple', got %q", cats[1].Slug)
	}
}

func TestParseCategoriesEmpty(t *testing.T) {
	html := `<html><body><div>no categories</div></body></html>`

	cats, err := ParseCategories(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cats) != 0 {
		t.Errorf("expected 0 categories, got %d", len(cats))
	}
}

func TestParseCategoriesDedup(t *testing.T) {
	html := `<html><body>
		<ul class="catalog_menu">
			<li>
				<ul class="catalog_menu_sub_second">
					<li><div class="bli_pos_second"><a href="/phones/">Телефоны</a></div></li>
					<li><div class="bli_pos_second"><a href="/phones/">Телефоны дубль</a></div></li>
				</ul>
			</li>
		</ul>
	</body></html>`

	cats, err := ParseCategories(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cats) != 1 {
		t.Errorf("expected 1 category after dedup, got %d", len(cats))
	}
}

func TestParseProducts(t *testing.T) {
	html := `<html><body>
		<div class="wrap_list_prod">
			<div class="blocks_product">
				<button class="favorite_product" data-elid="12345"></button>
				<div class="blocks_product_fix_w">
					<a href="/product/test-phone/" onclick="YandexEcommerce.getInstance().click([{&quot;name&quot;:&quot;Test Phone&quot;,&quot;id&quot;:&quot;SKU123&quot;,&quot;price&quot;:50000,&quot;brand&quot;:&quot;TestBrand&quot;,&quot;category&quot;:&quot;Phones&quot;}]);">
						<div class="bp_product_img">
							<img src="/img/phone.jpg" title="Test Phone" alt="Test Phone">
						</div>
					</a>
					<div class="bp_text">
						<p class="bp_text_price">50 000 —</p>
						<h2 class="bp_text_info"><a href="/product/test-phone/">Test Phone 128Gb</a></h2>
					</div>
				</div>
			</div>
		</div>
	</body></html>`

	products, err := ParseProducts(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	p := products[0]
	if p.ExternalID != "12345" {
		t.Errorf("expected ExternalID '12345', got %q", p.ExternalID)
	}
	if p.Name != "Test Phone 128Gb" {
		t.Errorf("expected Name 'Test Phone 128Gb', got %q", p.Name)
	}
	if p.SKU != "SKU123" {
		t.Errorf("expected SKU 'SKU123', got %q", p.SKU)
	}
	if p.Brand != "TestBrand" {
		t.Errorf("expected Brand 'TestBrand', got %q", p.Brand)
	}
	if p.Price != 50000 {
		t.Errorf("expected Price 50000, got %d", p.Price)
	}
	if p.ImageURL != "/img/phone.jpg" {
		t.Errorf("expected ImageURL '/img/phone.jpg', got %q", p.ImageURL)
	}
	if p.ProductURL != "/product/test-phone/" {
		t.Errorf("expected ProductURL '/product/test-phone/', got %q", p.ProductURL)
	}
}

func TestParseProductsEmpty(t *testing.T) {
	html := `<html><body><div class="wrap_list_prod"></div></body></html>`

	products, err := ParseProducts(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestParseProductsSkipNoURL(t *testing.T) {
	html := `<html><body>
		<div class="wrap_list_prod">
			<div class="blocks_product">
				<div class="blocks_product_fix_w">
					<div class="bp_text">
						<h2 class="bp_text_info"><a>No URL product</a></h2>
					</div>
				</div>
			</div>
		</div>
	</body></html>`

	products, err := ParseProducts(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 0 {
		t.Errorf("expected 0 products (no link), got %d", len(products))
	}
}

func TestParsePagination(t *testing.T) {
	html := `<html><body>
		<div class="pagination_catalog">
			<ul class="pagination">
				<li><a href="?p=1">1</a></li>
				<li><a href="?p=2" class="active">2</a></li>
				<li><a href="?p=3">3</a></li>
				<li><a href="?p=10">10</a></li>
			</ul>
		</div>
	</body></html>`

	info, err := ParsePagination(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.CurrentPage != 2 {
		t.Errorf("expected CurrentPage 2, got %d", info.CurrentPage)
	}
	if info.TotalPages != 10 {
		t.Errorf("expected TotalPages 10, got %d", info.TotalPages)
	}
}

func TestParsePaginationNone(t *testing.T) {
	html := `<html><body><div>no pagination</div></body></html>`

	info, err := ParsePagination(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.CurrentPage != 1 {
		t.Errorf("expected CurrentPage 1, got %d", info.CurrentPage)
	}
	if info.TotalPages != 1 {
		t.Errorf("expected TotalPages 1, got %d", info.TotalPages)
	}
}

func TestParseTotalProducts(t *testing.T) {
	html := `<html><body>
		<div class="search_all_produkt"><span>1 234</span></div>
	</body></html>`

	total, err := ParseTotalProducts(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 1234 {
		t.Errorf("expected 1234, got %d", total)
	}
}

func TestParseTotalProductsEmpty(t *testing.T) {
	html := `<html><body></body></html>`

	total, err := ParseTotalProducts(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 0 {
		t.Errorf("expected 0, got %d", total)
	}
}

func TestParseProductDescription(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "detail_text div",
			html:     `<html><body><div class="detail_text">Описание товара с подробностями</div></body></html>`,
			expected: "Описание товара с подробностями",
		},
		{
			name:     "itemprop description",
			html:     `<html><body><div itemprop="description">Характеристики устройства</div></body></html>`,
			expected: "Характеристики устройства",
		},
		{
			name:     "no description",
			html:     `<html><body><div>Нет описания</div></body></html>`,
			expected: "",
		},
		{
			name:     "whitespace normalization",
			html:     `<html><body><div class="detail_text">  Много    пробелов   </div></body></html>`,
			expected: "Много пробелов",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseProductDescription(tt.html)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParsePrice(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"50 000 —", 50000},
		{"1234", 1234},
		{"  100\u00a0500 — ", 100500},
		{"", 0},
		{"abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parsePrice(tt.input)
			if result != tt.expected {
				t.Errorf("parsePrice(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEcomPrice(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int
	}{
		{"float64", float64(39980), 39980},
		{"string", "40 780", 40780},
		{"nil", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ecomPrice(tt.input)
			if result != tt.expected {
				t.Errorf("ecomPrice(%v) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello   world  ", "hello world"},
		{"no extra spaces", "no extra spaces"},
		{"\t\nnewlines\t\n", "newlines"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeSpace(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeSpace(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
