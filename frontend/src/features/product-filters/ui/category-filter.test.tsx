import { describe, it, expect, vi } from "vitest";
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { renderWithI18n } from "@/test/render";
import { CategoryFilter } from "./category-filter";

const categories = [
  { id: "1", name: "Phones", slug: "phones", url: "/phones/", created_at: "", updated_at: "" },
  { id: "2", name: "Laptops", slug: "laptops", url: "/laptops/", created_at: "", updated_at: "" },
];

describe("CategoryFilter", () => {
  it("renders 'All categories' option", () => {
    renderWithI18n(<CategoryFilter categories={[]} value="" onChange={() => {}} />);
    expect(screen.getByText("All categories")).toBeInTheDocument();
  });

  it("renders category options", () => {
    renderWithI18n(<CategoryFilter categories={categories} value="" onChange={() => {}} />);
    expect(screen.getByText("Phones")).toBeInTheDocument();
    expect(screen.getByText("Laptops")).toBeInTheDocument();
  });

  it("selects provided value", () => {
    renderWithI18n(<CategoryFilter categories={categories} value="2" onChange={() => {}} />);
    const select = screen.getByRole("combobox") as HTMLSelectElement;
    expect(select.value).toBe("2");
  });

  it("calls onChange on selection", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    renderWithI18n(<CategoryFilter categories={categories} value="" onChange={onChange} />);

    await user.selectOptions(screen.getByRole("combobox"), "1");

    expect(onChange).toHaveBeenCalledWith("1");
  });

  it("renders label", () => {
    renderWithI18n(<CategoryFilter categories={[]} value="" onChange={() => {}} />);
    expect(screen.getByText("Category")).toBeInTheDocument();
  });
});
