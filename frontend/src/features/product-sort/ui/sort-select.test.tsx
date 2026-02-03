import { describe, it, expect, vi } from "vitest";
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { renderWithI18n } from "@/test/render";
import { SortSelect } from "./sort-select";

describe("SortSelect", () => {
  it("renders all sort options", () => {
    renderWithI18n(<SortSelect value="" onChange={() => {}} />);
    expect(screen.getByText("Newest")).toBeInTheDocument();
    expect(screen.getByText("Price: Low to High")).toBeInTheDocument();
    expect(screen.getByText("Price: High to Low")).toBeInTheDocument();
    expect(screen.getByText("Name: A-Z")).toBeInTheDocument();
    expect(screen.getByText("Name: Z-A")).toBeInTheDocument();
    expect(screen.getByText("Brand")).toBeInTheDocument();
  });

  it("defaults to 'Newest' when value is empty", () => {
    renderWithI18n(<SortSelect value="" onChange={() => {}} />);
    const select = screen.getByRole("combobox") as HTMLSelectElement;
    expect(select.value).toBe("created_at:desc");
  });

  it("selects provided value", () => {
    renderWithI18n(<SortSelect value="price:asc" onChange={() => {}} />);
    const select = screen.getByRole("combobox") as HTMLSelectElement;
    expect(select.value).toBe("price:asc");
  });

  it("calls onChange when selection changes", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    renderWithI18n(<SortSelect value="created_at:desc" onChange={onChange} />);

    await user.selectOptions(screen.getByRole("combobox"), "price:desc");

    expect(onChange).toHaveBeenCalledWith("price:desc");
  });
});
