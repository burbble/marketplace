import { describe, it, expect } from "vitest";
import { screen } from "@testing-library/react";
import { renderWithI18n } from "@/test/render";
import { ProductCard } from "./product-card";
import type { Product } from "../model";

vi.mock("next/link", () => ({
  default: ({ children, href, ...props }: { children: React.ReactNode; href: string }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}));

const baseProduct: Product = {
  id: "abc-123",
  external_id: "ext-1",
  sku: "SKU001",
  name: "iPhone 16 Pro",
  original_price: 150000,
  price: 149000,
  image_url: "/upload/iphone.jpg",
  product_url: "/iphone-16-pro/",
  brand: "Apple",
  description: "A great phone",
  category_id: "cat-1",
  created_at: "2025-01-01T00:00:00Z",
  updated_at: "2025-01-01T00:00:00Z",
};

describe("ProductCard", () => {
  it("renders product name", () => {
    renderWithI18n(<ProductCard product={baseProduct} exchangeRate={95.4} />);
    expect(screen.getByText("iPhone 16 Pro")).toBeInTheDocument();
  });

  it("renders product brand", () => {
    renderWithI18n(<ProductCard product={baseProduct} exchangeRate={95.4} />);
    expect(screen.getByText("Apple")).toBeInTheDocument();
  });

  it("links to product page", () => {
    renderWithI18n(<ProductCard product={baseProduct} exchangeRate={95.4} />);
    const link = screen.getByRole("link");
    expect(link).toHaveAttribute("href", "/products/abc-123");
  });

  it("renders image with alt text", () => {
    renderWithI18n(<ProductCard product={baseProduct} exchangeRate={95.4} />);
    const img = screen.getByAltText("iPhone 16 Pro");
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute("src", "https://store77.net/upload/iphone.jpg");
  });

  it("shows 'No image' when image_url is empty", () => {
    const product = { ...baseProduct, image_url: "" };
    renderWithI18n(<ProductCard product={product} exchangeRate={95.4} />);
    expect(screen.getByText("No image")).toBeInTheDocument();
  });

  it("shows original price when higher than current price", () => {
    renderWithI18n(<ProductCard product={baseProduct} exchangeRate={95.4} />);
    const prices = screen.getAllByText(/₽/);
    expect(prices.length).toBeGreaterThanOrEqual(2);
  });

  it("does not show original price when prices are equal", () => {
    const product = { ...baseProduct, original_price: 149000 };
    renderWithI18n(<ProductCard product={product} exchangeRate={95.4} />);
    const prices = screen.getAllByText(/₽/);
    expect(prices).toHaveLength(1);
  });

  it("shows USDT price when exchange rate is positive", () => {
    renderWithI18n(<ProductCard product={baseProduct} exchangeRate={95.4} />);
    expect(screen.getByText(/USDT/)).toBeInTheDocument();
  });

  it("does not show USDT price when exchange rate is 0", () => {
    renderWithI18n(<ProductCard product={baseProduct} exchangeRate={0} />);
    expect(screen.queryByText(/USDT/)).not.toBeInTheDocument();
  });
});
