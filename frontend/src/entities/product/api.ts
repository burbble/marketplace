import { apiFetch } from "@/shared/api/client";
import type { Product, ProductFilter, ProductList } from "./model";

export async function getProducts(filter: ProductFilter): Promise<ProductList> {
  const params = new URLSearchParams();

  if (filter.page) params.set("page", String(filter.page));
  if (filter.page_size) params.set("page_size", String(filter.page_size));
  if (filter.sort_fields) params.set("sort_fields", filter.sort_fields);
  if (filter.category_id) params.set("category_id", filter.category_id);
  if (filter.brand) params.set("brand", filter.brand);
  if (filter.min_price != null) params.set("min_price", String(filter.min_price));
  if (filter.max_price != null) params.set("max_price", String(filter.max_price));
  if (filter.search) params.set("search", filter.search);

  const qs = params.toString();
  const data = await apiFetch<ProductList>(`/products${qs ? `?${qs}` : ""}`);
  return { ...data, products: data.products ?? [] };
}

export async function getProductById(id: string): Promise<Product> {
  return apiFetch<Product>(`/products/${id}`);
}

export async function getBrands(): Promise<string[]> {
  const data = await apiFetch<string[]>("/brands");
  return data ?? [];
}
