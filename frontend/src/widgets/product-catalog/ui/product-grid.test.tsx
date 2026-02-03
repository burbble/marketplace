import { describe, it, expect } from "vitest";
import { screen } from "@testing-library/react";
import { renderWithI18n } from "@/test/render";
import { ProductGrid } from "./product-grid";
import type { Product } from "@/entities/product/model";

vi.mock("next/link", () => ({
  default: ({ children, href, ...props }: { children: React.ReactNode; href: string }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}));

const makeProduct = (id: string, name: string): Product => ({
  id,
  external_id: `ext-${id}`,
  sku: `SKU-${id}`,
  name,
  original_price: 10000,
  price: 9000,
  image_url: `/img/${id}.jpg`,
  product_url: `/${id}/`,
  brand: "Apple",
  description: "",
  category_id: "cat-1",
  created_at: "2025-01-01T00:00:00Z",
  updated_at: "2025-01-01T00:00:00Z",
});

describe("ProductGrid", () => {
  it("renders empty state when products is empty", () => {
    renderWithI18n(<ProductGrid products={[]} exchangeRate={95} />);
    expect(screen.getByText("No products found")).toBeInTheDocument();
    expect(screen.getByText("Try adjusting your filters")).toBeInTheDocument();
  });

  it("renders product cards", () => {
    const products = [makeProduct("1", "iPhone"), makeProduct("2", "MacBook")];
    renderWithI18n(<ProductGrid products={products} exchangeRate={95} />);
    expect(screen.getByText("iPhone")).toBeInTheDocument();
    expect(screen.getByText("MacBook")).toBeInTheDocument();
  });

  it("renders correct number of links", () => {
    const products = [makeProduct("1", "A"), makeProduct("2", "B"), makeProduct("3", "C")];
    renderWithI18n(<ProductGrid products={products} exchangeRate={95} />);
    const links = screen.getAllByRole("link");
    expect(links).toHaveLength(3);
  });
});
