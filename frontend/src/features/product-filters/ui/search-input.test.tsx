import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { screen, fireEvent } from "@testing-library/react";
import { renderWithI18n } from "@/test/render";
import { SearchInput } from "./search-input";

describe("SearchInput", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("renders input with placeholder", () => {
    renderWithI18n(<SearchInput value="" onChange={() => {}} />);
    expect(screen.getByPlaceholderText("Search products...")).toBeInTheDocument();
  });

  it("sets default value", () => {
    renderWithI18n(<SearchInput value="iphone" onChange={() => {}} />);
    const input = screen.getByPlaceholderText("Search products...") as HTMLInputElement;
    expect(input.defaultValue).toBe("iphone");
  });

  it("calls onChange after debounce", () => {
    const onChange = vi.fn();
    renderWithI18n(<SearchInput value="" onChange={onChange} />);

    fireEvent.change(screen.getByPlaceholderText("Search products..."), {
      target: { value: "test" },
    });

    expect(onChange).not.toHaveBeenCalled();

    vi.advanceTimersByTime(400);

    expect(onChange).toHaveBeenCalledWith("test");
  });

  it("debounces rapid input", () => {
    const onChange = vi.fn();
    renderWithI18n(<SearchInput value="" onChange={onChange} />);
    const input = screen.getByPlaceholderText("Search products...");

    fireEvent.change(input, { target: { value: "a" } });
    vi.advanceTimersByTime(200);

    fireEvent.change(input, { target: { value: "ab" } });
    vi.advanceTimersByTime(400);

    expect(onChange).toHaveBeenCalledTimes(1);
    expect(onChange).toHaveBeenCalledWith("ab");
  });
});
