import { apiFetch } from "@/shared/api/client";
import type { ExchangeRate } from "./model";

export async function getExchangeRate(): Promise<ExchangeRate> {
  return apiFetch<ExchangeRate>("/exchange/rate");
}
