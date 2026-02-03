"use client";

import { useTranslation } from "@/shared/i18n";

export function BrandFilter({
  brands,
  value,
  onChange,
}: {
  brands: string[];
  value: string;
  onChange: (v: string) => void;
}) {
  const { t } = useTranslation();

  return (
    <div className="flex flex-col gap-1.5">
      <label className="text-xs font-medium text-zinc-400">{t("filter.brand")}</label>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="w-full rounded-lg border border-zinc-700 bg-zinc-800 py-2.5 pl-3 pr-9 text-sm text-zinc-200 outline-none transition-colors focus:border-zinc-500"
      >
        <option value="">{t("filter.allBrands")}</option>
        {brands.map((b) => (
          <option key={b} value={b}>
            {b}
          </option>
        ))}
      </select>
    </div>
  );
}
