import { describe, it, expect, vi, beforeEach } from "vitest";
import { getExchangeRate } from "./api";

vi.mock("@/shared/api/client", () => ({
  apiFetch: vi.fn(),
}));

import { apiFetch } from "@/shared/api/client";

const mockApiFetch = vi.mocked(apiFetch);

describe("getExchangeRate", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calls /exchange/rate", async () => {
    mockApiFetch.mockResolvedValue({ rate: 95.4 });

    const result = await getExchangeRate();

    expect(mockApiFetch).toHaveBeenCalledWith("/exchange/rate");
    expect(result).toEqual({ rate: 95.4 });
  });
});
