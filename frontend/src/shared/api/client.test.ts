import { describe, it, expect, vi, beforeEach } from "vitest";
import { apiFetch } from "./client";

describe("apiFetch", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("fetches data successfully", async () => {
    const mockData = { id: 1, name: "Test" };
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockData),
      })
    );

    const result = await apiFetch<typeof mockData>("/test");

    expect(result).toEqual(mockData);
    expect(fetch).toHaveBeenCalledWith("/api/v1/test", {
      headers: { "Content-Type": "application/json" },
    });
  });

  it("sets Content-Type header", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({}),
      })
    );

    await apiFetch("/test");

    const call = vi.mocked(fetch).mock.calls[0];
    expect((call[1]?.headers as Record<string, string>)["Content-Type"]).toBe("application/json");
  });

  it("merges custom headers", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({}),
      })
    );

    await apiFetch("/test", {
      headers: { Authorization: "Bearer token" },
    });

    const call = vi.mocked(fetch).mock.calls[0];
    const headers = call[1]?.headers as Record<string, string>;
    expect(headers["Content-Type"]).toBe("application/json");
    expect(headers["Authorization"]).toBe("Bearer token");
  });

  it("throws error with API error message", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 404,
        json: () => Promise.resolve({ error: "not found" }),
      })
    );

    await expect(apiFetch("/test")).rejects.toThrow("not found");
  });

  it("throws generic error when body has no error field", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: () => Promise.resolve({}),
      })
    );

    await expect(apiFetch("/test")).rejects.toThrow("API error: 500");
  });

  it("throws generic error when body parsing fails", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: () => Promise.reject(new Error("parse error")),
      })
    );

    await expect(apiFetch("/test")).rejects.toThrow("API error: 500");
  });
});
