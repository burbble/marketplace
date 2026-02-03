import { describe, it, expect } from "vitest";
import { formatRUB, formatUSDT, cn, resolveImageUrl } from "./format";

describe("formatRUB", () => {
  it("formats price as Russian Rubles", () => {
    const result = formatRUB(150000);
    expect(result).toContain("150");
    expect(result).toContain("000");
  });

  it("formats zero", () => {
    const result = formatRUB(0);
    expect(result).toContain("0");
  });

  it("has no fraction digits", () => {
    const result = formatRUB(1234.56);
    expect(result).not.toContain(".");
    expect(result).not.toContain(",56");
  });
});

describe("formatUSDT", () => {
  it("converts RUB to USDT", () => {
    expect(formatUSDT(9540, 95.4)).toBe("100.00 USDT");
  });

  it("returns dash when rate is zero", () => {
    expect(formatUSDT(1000, 0)).toBe("—");
  });

  it("returns dash when rate is negative", () => {
    expect(formatUSDT(1000, -1)).toBe("—");
  });

  it("formats with two decimal places", () => {
    expect(formatUSDT(10000, 95)).toBe("105.26 USDT");
  });
});

describe("cn", () => {
  it("joins class names", () => {
    expect(cn("a", "b", "c")).toBe("a b c");
  });

  it("filters out falsy values", () => {
    expect(cn("a", false, undefined, null, "b")).toBe("a b");
  });

  it("returns empty string for no classes", () => {
    expect(cn()).toBe("");
  });

  it("returns empty string for all falsy", () => {
    expect(cn(false, undefined, null)).toBe("");
  });
});

describe("resolveImageUrl", () => {
  it("returns absolute URLs as-is", () => {
    expect(resolveImageUrl("https://example.com/img.jpg")).toBe("https://example.com/img.jpg");
  });

  it("returns http URLs as-is", () => {
    expect(resolveImageUrl("http://example.com/img.jpg")).toBe("http://example.com/img.jpg");
  });

  it("prepends store77 domain for relative paths", () => {
    expect(resolveImageUrl("/upload/img.jpg")).toBe("https://store77.net/upload/img.jpg");
  });

  it("returns empty string for empty input", () => {
    expect(resolveImageUrl("")).toBe("");
  });
});
