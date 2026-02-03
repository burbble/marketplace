export const dictionaries = {
  en: {
    // Layout
    "meta.title": "Marketplace",
    "meta.description": "Product catalog",

    // Header
    "header.title": "Marketplace",

    // Catalog page
    "catalog.filters": "Filters",
    "catalog.products_zero": "0 products",
    "catalog.products_one": "1 product",
    "catalog.products_other": "{count} products",

    // Filters
    "filter.search": "Search products...",
    "filter.category": "Category",
    "filter.allCategories": "All categories",
    "filter.brand": "Brand",
    "filter.allBrands": "All brands",
    "filter.price": "Price (RUB)",
    "filter.priceFrom": "From",
    "filter.priceTo": "To",
    "filter.reset": "Reset filters",

    // Sort
    "sort.newest": "Newest",
    "sort.priceLow": "Price: Low to High",
    "sort.priceHigh": "Price: High to Low",
    "sort.nameAZ": "Name: A-Z",
    "sort.nameZA": "Name: Z-A",
    "sort.brand": "Brand",

    // Product grid
    "grid.empty": "No products found",
    "grid.emptyHint": "Try adjusting your filters",

    // Product card / detail
    "product.noImage": "No image",
    "product.notFound": "Product not found",
    "product.loadError": "Failed to load product",
    "product.backToCatalog": "Back to catalog",
    "product.description": "Description",
    "product.details": "Details",
    "product.brand": "Brand",
    "product.sku": "SKU",
    "product.externalId": "External ID",

    // Exchange
    "exchange.unavailable": "USDT —",
    "exchange.rate": "1 USDT = {rate} RUB",

    // Spinner
    loading: "Loading...",
  },

  ru: {
    // Layout
    "meta.title": "Маркетплейс",
    "meta.description": "Каталог товаров",

    // Header
    "header.title": "Маркетплейс",

    // Catalog page
    "catalog.filters": "Фильтры",
    "catalog.products_zero": "0 товаров",
    "catalog.products_one": "1 товар",
    "catalog.products_other": "{count} товаров",

    // Filters
    "filter.search": "Поиск товаров...",
    "filter.category": "Категория",
    "filter.allCategories": "Все категории",
    "filter.brand": "Бренд",
    "filter.allBrands": "Все бренды",
    "filter.price": "Цена (RUB)",
    "filter.priceFrom": "От",
    "filter.priceTo": "До",
    "filter.reset": "Сбросить фильтры",

    // Sort
    "sort.newest": "Новые",
    "sort.priceLow": "Цена: по возрастанию",
    "sort.priceHigh": "Цена: по убыванию",
    "sort.nameAZ": "Название: А-Я",
    "sort.nameZA": "Название: Я-А",
    "sort.brand": "Бренд",

    // Product grid
    "grid.empty": "Товары не найдены",
    "grid.emptyHint": "Попробуйте изменить фильтры",

    // Product card / detail
    "product.noImage": "Нет изображения",
    "product.notFound": "Товар не найден",
    "product.loadError": "Не удалось загрузить товар",
    "product.backToCatalog": "Назад в каталог",
    "product.description": "Описание",
    "product.details": "Подробности",
    "product.brand": "Бренд",
    "product.sku": "Артикул",
    "product.externalId": "Внешний ID",

    // Exchange
    "exchange.unavailable": "USDT —",
    "exchange.rate": "1 USDT = {rate} RUB",

    // Spinner
    loading: "Загрузка...",
  },
} as const;

export type Locale = keyof typeof dictionaries;
export type TranslationKey = keyof (typeof dictionaries)["en"];
