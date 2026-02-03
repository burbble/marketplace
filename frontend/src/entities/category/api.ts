import { apiFetch } from "@/shared/api/client";
import type { Category } from "./model";

export async function getCategories(): Promise<Category[]> {
  const data = await apiFetch<Category[]>("/categories");
  return data ?? [];
}
