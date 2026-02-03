"use client";

import { useTranslation } from "@/shared/i18n";
import type { Locale } from "@/shared/i18n";

const LOCALES: { value: Locale; label: string }[] = [
  { value: "en", label: "EN" },
  { value: "ru", label: "RU" },
];

export function LanguageSwitcher() {
  const { locale, setLocale } = useTranslation();

  return (
    <div className="flex overflow-hidden rounded-lg border border-zinc-700">
      {LOCALES.map((l) => (
        <button
          key={l.value}
          onClick={() => setLocale(l.value)}
          className={`px-2.5 py-1 text-xs font-medium transition-colors ${
            locale === l.value
              ? "bg-zinc-700 text-white"
              : "bg-zinc-800 text-zinc-400 hover:text-zinc-200"
          }`}
        >
          {l.label}
        </button>
      ))}
    </div>
  );
}
