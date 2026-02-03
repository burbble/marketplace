import { describe, it, expect, vi, beforeEach } from "vitest";
import { screen, waitFor } from "@testing-library/react";
import { renderWithI18n } from "@/test/render";
import { ExchangeBadge } from "./exchange-badge";

vi.mock("@/entities/exchange/api", () => ({
  getExchangeRate: vi.fn(),
}));

import { getExchangeRate } from "@/entities/exchange/api";

const mockGetRate = vi.mocked(getExchangeRate);

describe("ExchangeBadge", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers();
  });

  it("shows loading state initially", () => {
    mockGetRate.mockReturnValue(new Promise(() => {}));
    renderWithI18n(<ExchangeBadge />);
    expect(screen.getByText(/USDT —/)).toBeInTheDocument();
  });

  it("shows rate after fetch", async () => {
    mockGetRate.mockResolvedValue({ rate: 95.4 });
    vi.useRealTimers();

    renderWithI18n(<ExchangeBadge />);

    await waitFor(() => {
      expect(screen.getByText(/1 USDT = 95\.40 RUB/)).toBeInTheDocument();
    });
  });

  it("shows warning badge on fetch error", async () => {
    mockGetRate.mockRejectedValue(new Error("network error"));
    vi.useRealTimers();

    renderWithI18n(<ExchangeBadge />);

    await waitFor(() => {
      expect(screen.getByText(/USDT —/)).toBeInTheDocument();
    });
  });

  it("refreshes rate on interval", async () => {
    mockGetRate.mockResolvedValue({ rate: 95.4 });

    renderWithI18n(<ExchangeBadge />);

    await vi.advanceTimersByTimeAsync(0);

    expect(mockGetRate).toHaveBeenCalledTimes(1);

    await vi.advanceTimersByTimeAsync(60_000);

    expect(mockGetRate).toHaveBeenCalledTimes(2);
  });
});
