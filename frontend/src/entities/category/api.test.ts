import { describe, it, expect, vi, beforeEach } from "vitest";
import { getCategories } from "./api";

vi.mock("@/shared/api/client", () => ({
  apiFetch: vi.fn(),
}));

import { apiFetch } from "@/shared/api/client";

const mockApiFetch = vi.mocked(apiFetch);

describe("getCategories", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calls /categories", async () => {
    const categories = [
      { id: "1", name: "Phones", slug: "phones" },
      { id: "2", name: "Laptops", slug: "laptops" },
    ];
    mockApiFetch.mockResolvedValue(categories);

    const result = await getCategories();

    expect(mockApiFetch).toHaveBeenCalledWith("/categories");
    expect(result).toEqual(categories);
  });
});
