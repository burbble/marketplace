"use client";

import { useState, useEffect } from "react";
import type { Category } from "@/entities/category/model";
import { getCategories } from "@/entities/category/api";
import { getBrands } from "@/entities/product/api";
import { SearchInput } from "./search-input";
import { CategoryFilter } from "./category-filter";
import { BrandFilter } from "./brand-filter";
import { PriceFilter } from "./price-filter";

export function FilterSidebar({
  search,
  categoryId,
  brand,
  minPrice,
  maxPrice,
  onSearchChange,
  onCategoryChange,
  onBrandChange,
  onMinPriceChange,
  onMaxPriceChange,
  onReset,
}: {
  search: string;
  categoryId: string;
  brand: string;
  minPrice: string;
  maxPrice: string;
  onSearchChange: (v: string) => void;
  onCategoryChange: (v: string) => void;
  onBrandChange: (v: string) => void;
  onMinPriceChange: (v: string) => void;
  onMaxPriceChange: (v: string) => void;
  onReset: () => void;
}) {
  const [categories, setCategories] = useState<Category[]>([]);
  const [brands, setBrands] = useState<string[]>([]);

  useEffect(() => {
    getCategories()
      .then(setCategories)
      .catch(() => {});
    getBrands()
      .then(setBrands)
      .catch(() => {});
  }, []);

  const hasFilters = search || categoryId || brand || minPrice || maxPrice;

  return (
    <div className="flex flex-col gap-5">
      <SearchInput value={search} onChange={onSearchChange} />
      <CategoryFilter
        categories={categories}
        value={categoryId}
        onChange={onCategoryChange}
      />
      <BrandFilter brands={brands} value={brand} onChange={onBrandChange} />
      <PriceFilter
        minPrice={minPrice}
        maxPrice={maxPrice}
        onMinChange={onMinPriceChange}
        onMaxChange={onMaxPriceChange}
      />
      {hasFilters && (
        <button
          onClick={onReset}
          className="rounded-lg border border-zinc-700 bg-zinc-800 px-4 py-2.5 text-sm text-zinc-400 transition-colors hover:border-zinc-600 hover:text-zinc-200"
        >
          Reset filters
        </button>
      )}
    </div>
  );
}
