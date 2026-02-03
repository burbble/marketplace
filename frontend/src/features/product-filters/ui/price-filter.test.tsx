import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { screen, fireEvent } from "@testing-library/react";
import { renderWithI18n } from "@/test/render";
import { PriceFilter } from "./price-filter";

describe("PriceFilter", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("renders label", () => {
    renderWithI18n(
      <PriceFilter minPrice="" maxPrice="" onMinChange={() => {}} onMaxChange={() => {}} />
    );
    expect(screen.getByText("Price (RUB)")).toBeInTheDocument();
  });

  it("renders min and max inputs", () => {
    renderWithI18n(
      <PriceFilter minPrice="" maxPrice="" onMinChange={() => {}} onMaxChange={() => {}} />
    );
    expect(screen.getByPlaceholderText("From")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("To")).toBeInTheDocument();
  });

  it("sets default values", () => {
    renderWithI18n(
      <PriceFilter minPrice="1000" maxPrice="5000" onMinChange={() => {}} onMaxChange={() => {}} />
    );
    const from = screen.getByPlaceholderText("From") as HTMLInputElement;
    const to = screen.getByPlaceholderText("To") as HTMLInputElement;
    expect(from.defaultValue).toBe("1000");
    expect(to.defaultValue).toBe("5000");
  });

  it("calls onMinChange after debounce", () => {
    const onMinChange = vi.fn();
    renderWithI18n(
      <PriceFilter minPrice="" maxPrice="" onMinChange={onMinChange} onMaxChange={() => {}} />
    );

    fireEvent.change(screen.getByPlaceholderText("From"), {
      target: { value: "500" },
    });

    expect(onMinChange).not.toHaveBeenCalled();

    vi.advanceTimersByTime(500);

    expect(onMinChange).toHaveBeenCalledWith("500");
  });

  it("calls onMaxChange after debounce", () => {
    const onMaxChange = vi.fn();
    renderWithI18n(
      <PriceFilter minPrice="" maxPrice="" onMinChange={() => {}} onMaxChange={onMaxChange} />
    );

    fireEvent.change(screen.getByPlaceholderText("To"), {
      target: { value: "9999" },
    });

    expect(onMaxChange).not.toHaveBeenCalled();

    vi.advanceTimersByTime(500);

    expect(onMaxChange).toHaveBeenCalledWith("9999");
  });
});
