"use client";

import { useRef } from "react";
import { useTranslation } from "@/shared/i18n";

export function PriceFilter({
  minPrice,
  maxPrice,
  onMinChange,
  onMaxChange,
}: {
  minPrice: string;
  maxPrice: string;
  onMinChange: (v: string) => void;
  onMaxChange: (v: string) => void;
}) {
  const { t } = useTranslation();
  const minTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const maxTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  function handleMin(e: React.ChangeEvent<HTMLInputElement>) {
    const val = e.target.value;
    if (minTimer.current) clearTimeout(minTimer.current);
    minTimer.current = setTimeout(() => onMinChange(val), 500);
  }

  function handleMax(e: React.ChangeEvent<HTMLInputElement>) {
    const val = e.target.value;
    if (maxTimer.current) clearTimeout(maxTimer.current);
    maxTimer.current = setTimeout(() => onMaxChange(val), 500);
  }

  return (
    <div className="flex flex-col gap-1.5">
      <label className="text-xs font-medium text-zinc-400">{t("filter.price")}</label>
      <div className="flex items-center gap-2">
        <input
          type="number"
          defaultValue={minPrice}
          onChange={handleMin}
          placeholder={t("filter.priceFrom")}
          min={0}
          className="w-full rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2.5 text-sm text-zinc-200 placeholder-zinc-500 outline-none transition-colors focus:border-zinc-500"
        />
        <span className="text-zinc-600">—</span>
        <input
          type="number"
          defaultValue={maxPrice}
          onChange={handleMax}
          placeholder={t("filter.priceTo")}
          min={0}
          className="w-full rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2.5 text-sm text-zinc-200 placeholder-zinc-500 outline-none transition-colors focus:border-zinc-500"
        />
      </div>
    </div>
  );
}
