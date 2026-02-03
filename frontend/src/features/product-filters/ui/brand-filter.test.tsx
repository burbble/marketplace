import { describe, it, expect, vi } from "vitest";
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { renderWithI18n } from "@/test/render";
import { BrandFilter } from "./brand-filter";

describe("BrandFilter", () => {
  it("renders 'All brands' option", () => {
    renderWithI18n(<BrandFilter brands={[]} value="" onChange={() => {}} />);
    expect(screen.getByText("All brands")).toBeInTheDocument();
  });

  it("renders brand options", () => {
    renderWithI18n(<BrandFilter brands={["Apple", "Samsung"]} value="" onChange={() => {}} />);
    expect(screen.getByText("Apple")).toBeInTheDocument();
    expect(screen.getByText("Samsung")).toBeInTheDocument();
  });

  it("selects provided value", () => {
    renderWithI18n(
      <BrandFilter brands={["Apple", "Samsung"]} value="Samsung" onChange={() => {}} />
    );
    const select = screen.getByRole("combobox") as HTMLSelectElement;
    expect(select.value).toBe("Samsung");
  });

  it("calls onChange on selection", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    renderWithI18n(<BrandFilter brands={["Apple", "Samsung"]} value="" onChange={onChange} />);

    await user.selectOptions(screen.getByRole("combobox"), "Apple");

    expect(onChange).toHaveBeenCalledWith("Apple");
  });

  it("renders label", () => {
    renderWithI18n(<BrandFilter brands={[]} value="" onChange={() => {}} />);
    expect(screen.getByText("Brand")).toBeInTheDocument();
  });
});
