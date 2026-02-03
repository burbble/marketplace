"use client";

import Link from "next/link";
import { ExchangeBadge } from "@/features/exchange-indicator/ui/exchange-badge";
import { LanguageSwitcher } from "@/features/language-switcher/ui/language-switcher";
import { useTranslation } from "@/shared/i18n";

export function Header() {
  const { t } = useTranslation();

  return (
    <header className="sticky top-0 z-50 border-b border-zinc-800 bg-zinc-950/80 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6">
        <Link href="/" className="text-lg font-bold text-white">
          {t("header.title")}
        </Link>
        <div className="flex items-center gap-3">
          <ExchangeBadge />
          <LanguageSwitcher />
        </div>
      </div>
    </header>
  );
}
