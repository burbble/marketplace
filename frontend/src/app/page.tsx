"use client";

import { useState, useEffect, useCallback, Suspense } from "react";
import { useQueryState, parseAsString, parseAsInteger } from "nuqs";
import { getProducts } from "@/entities/product/api";
import { getExchangeRate } from "@/entities/exchange/api";
import type { Product } from "@/entities/product/model";
import { FilterSidebar } from "@/features/product-filters/ui/filter-sidebar";
import { SortSelect } from "@/features/product-sort/ui/sort-select";
import { ProductGrid } from "@/widgets/product-catalog/ui/product-grid";
import { Pagination } from "@/widgets/product-catalog/ui/pagination";
import { Spinner } from "@/shared/ui/spinner";

function CatalogContent() {
  const [search, setSearch] = useQueryState("search", parseAsString.withDefault(""));
  const [categoryId, setCategoryId] = useQueryState("category_id", parseAsString.withDefault(""));
  const [brand, setBrand] = useQueryState("brand", parseAsString.withDefault(""));
  const [minPrice, setMinPrice] = useQueryState("min_price", parseAsString.withDefault(""));
  const [maxPrice, setMaxPrice] = useQueryState("max_price", parseAsString.withDefault(""));
  const [sortFields, setSortFields] = useQueryState("sort", parseAsString.withDefault(""));
  const [page, setPage] = useQueryState("page", parseAsInteger.withDefault(1));

  const [products, setProducts] = useState<Product[]>([]);
  const [total, setTotal] = useState(0);
  const [pageSize] = useState(24);
  const [loading, setLoading] = useState(true);
  const [exchangeRate, setExchangeRate] = useState(0);
  const [mobileFiltersOpen, setMobileFiltersOpen] = useState(false);

  const fetchProducts = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getProducts({
        page,
        page_size: pageSize,
        sort_fields: sortFields || undefined,
        category_id: categoryId || undefined,
        brand: brand || undefined,
        min_price: minPrice ? Number(minPrice) : undefined,
        max_price: maxPrice ? Number(maxPrice) : undefined,
        search: search || undefined,
      });
      setProducts(data.products || []);
      setTotal(data.total);
    } catch {
      setProducts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, sortFields, categoryId, brand, minPrice, maxPrice, search]);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  useEffect(() => {
    let active = true;
    async function load() {
      try {
        const data = await getExchangeRate();
        if (active) setExchangeRate(data.rate);
      } catch {}
    }
    load();
    const interval = setInterval(load, 60_000);
    return () => {
      active = false;
      clearInterval(interval);
    };
  }, []);

  function handleFilterChange(setter: (v: string) => Promise<URLSearchParams>) {
    return (v: string) => {
      setter(v);
      setPage(1);
    };
  }

  function handleReset() {
    setSearch("");
    setCategoryId("");
    setBrand("");
    setMinPrice("");
    setMaxPrice("");
    setPage(1);
  }

  const sidebar = (
    <FilterSidebar
      search={search}
      categoryId={categoryId}
      brand={brand}
      minPrice={minPrice}
      maxPrice={maxPrice}
      onSearchChange={handleFilterChange(setSearch)}
      onCategoryChange={handleFilterChange(setCategoryId)}
      onBrandChange={handleFilterChange(setBrand)}
      onMinPriceChange={handleFilterChange(setMinPrice)}
      onMaxPriceChange={handleFilterChange(setMaxPrice)}
      onReset={handleReset}
    />
  );

  return (
    <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6">
      <div className="lg:hidden mb-4">
        <button
          onClick={() => setMobileFiltersOpen(!mobileFiltersOpen)}
          className="flex w-full items-center justify-center gap-2 rounded-lg border border-zinc-700 bg-zinc-800 px-4 py-2.5 text-sm text-zinc-300 transition-colors hover:border-zinc-600"
        >
          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
          </svg>
          Filters
        </button>

        {mobileFiltersOpen && (
          <div className="mt-4 rounded-xl border border-zinc-800 bg-zinc-900 p-4">
            {sidebar}
          </div>
        )}
      </div>

      <div className="flex gap-8">
        <aside className="hidden w-64 shrink-0 lg:block">
          <div className="sticky top-20 rounded-xl border border-zinc-800 bg-zinc-900 p-5">
            {sidebar}
          </div>
        </aside>

        <div className="min-w-0 flex-1">
          <div className="mb-4 flex items-center justify-between">
            <p className="text-sm text-zinc-400">
              {loading ? (
                <Spinner className="h-4 w-4" />
              ) : (
                <>{total} product{total !== 1 ? "s" : ""}</>
              )}
            </p>
            <SortSelect
              value={sortFields}
              onChange={(v) => {
                setSortFields(v);
                setPage(1);
              }}
            />
          </div>

          {loading ? (
            <div className="flex items-center justify-center py-20">
              <Spinner className="h-8 w-8" />
            </div>
          ) : (
            <>
              <ProductGrid products={products} exchangeRate={exchangeRate} />
              <Pagination
                page={page}
                total={total}
                pageSize={pageSize}
                onChange={setPage}
              />
            </>
          )}
        </div>
      </div>
    </div>
  );
}

export default function CatalogPage() {
  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center py-20">
          <Spinner className="h-8 w-8" />
        </div>
      }
    >
      <CatalogContent />
    </Suspense>
  );
}
