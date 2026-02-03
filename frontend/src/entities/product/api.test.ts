import { describe, it, expect, vi, beforeEach } from "vitest";
import { getProducts, getProductById, getBrands } from "./api";

vi.mock("@/shared/api/client", () => ({
  apiFetch: vi.fn(),
}));

import { apiFetch } from "@/shared/api/client";

const mockApiFetch = vi.mocked(apiFetch);

describe("getProducts", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calls /products without params when filter is empty", async () => {
    mockApiFetch.mockResolvedValue({ products: [], total: 0, page: 1, page_size: 24 });

    await getProducts({});

    expect(mockApiFetch).toHaveBeenCalledWith("/products");
  });

  it("builds query string from filter", async () => {
    mockApiFetch.mockResolvedValue({ products: [], total: 0, page: 1, page_size: 24 });

    await getProducts({
      page: 2,
      page_size: 10,
      sort_fields: "price:asc",
      category_id: "abc-123",
      brand: "Apple",
      min_price: 1000,
      max_price: 5000,
      search: "iphone",
    });

    const url = mockApiFetch.mock.calls[0][0];
    expect(url).toContain("/products?");
    expect(url).toContain("page=2");
    expect(url).toContain("page_size=10");
    expect(url).toContain("sort_fields=price%3Aasc");
    expect(url).toContain("category_id=abc-123");
    expect(url).toContain("brand=Apple");
    expect(url).toContain("min_price=1000");
    expect(url).toContain("max_price=5000");
    expect(url).toContain("search=iphone");
  });

  it("defaults null products to empty array", async () => {
    mockApiFetch.mockResolvedValue({ products: null, total: 0, page: 1, page_size: 24 });

    const result = await getProducts({});

    expect(result.products).toEqual([]);
  });

  it("preserves non-null products", async () => {
    const products = [{ id: "1", name: "Phone" }];
    mockApiFetch.mockResolvedValue({ products, total: 1, page: 1, page_size: 24 });

    const result = await getProducts({});

    expect(result.products).toEqual(products);
  });
});

describe("getProductById", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calls /products/:id", async () => {
    const product = { id: "abc", name: "Phone" };
    mockApiFetch.mockResolvedValue(product);

    const result = await getProductById("abc");

    expect(mockApiFetch).toHaveBeenCalledWith("/products/abc");
    expect(result).toEqual(product);
  });
});

describe("getBrands", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calls /brands", async () => {
    const brands = ["Apple", "Samsung"];
    mockApiFetch.mockResolvedValue(brands);

    const result = await getBrands();

    expect(mockApiFetch).toHaveBeenCalledWith("/brands");
    expect(result).toEqual(brands);
  });
});
