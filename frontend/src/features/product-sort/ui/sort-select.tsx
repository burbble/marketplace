"use client";

import { useTranslation } from "@/shared/i18n";
import type { TranslationKey } from "@/shared/i18n";

const SORT_OPTIONS: { labelKey: TranslationKey; value: string }[] = [
  { labelKey: "sort.newest", value: "created_at:desc" },
  { labelKey: "sort.priceLow", value: "price:asc" },
  { labelKey: "sort.priceHigh", value: "price:desc" },
  { labelKey: "sort.nameAZ", value: "name:asc" },
  { labelKey: "sort.nameZA", value: "name:desc" },
  { labelKey: "sort.brand", value: "brand:asc" },
];

export function SortSelect({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  const { t } = useTranslation();

  return (
    <select
      value={value || "created_at:desc"}
      onChange={(e) => onChange(e.target.value)}
      className="rounded-lg border border-zinc-700 bg-zinc-800 py-2 pl-3 pr-9 text-sm text-zinc-200 outline-none transition-colors focus:border-zinc-500"
    >
      {SORT_OPTIONS.map((opt) => (
        <option key={opt.value} value={opt.value}>
          {t(opt.labelKey)}
        </option>
      ))}
    </select>
  );
}
