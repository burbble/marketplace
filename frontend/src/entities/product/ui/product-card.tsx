"use client";

import Link from "next/link";
import type { Product } from "../model";
import { formatRUB, formatUSDT, resolveImageUrl } from "@/shared/lib/format";
import { useTranslation } from "@/shared/i18n";

export function ProductCard({ product, exchangeRate }: { product: Product; exchangeRate: number }) {
  const { t } = useTranslation();

  return (
    <Link
      href={`/products/${product.id}`}
      className="group flex flex-col overflow-hidden rounded-xl border border-zinc-800 bg-zinc-900 transition-colors hover:border-zinc-700"
    >
      <div className="relative aspect-square overflow-hidden bg-zinc-800">
        {product.image_url ? (
          <img
            src={resolveImageUrl(product.image_url)}
            alt={product.name}
            loading="lazy"
            className="h-full w-full object-cover transition-transform group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full items-center justify-center text-zinc-600">
            {t("product.noImage")}
          </div>
        )}
      </div>

      <div className="flex flex-1 flex-col gap-2 p-4">
        <p className="text-xs text-zinc-500">{product.brand}</p>
        <h3 className="line-clamp-2 text-sm font-medium text-zinc-200">{product.name}</h3>

        <div className="mt-auto flex flex-col gap-1">
          <div className="flex items-baseline gap-2">
            <span className="text-lg font-bold text-white">{formatRUB(product.price)}</span>
            {product.original_price > product.price && (
              <span className="text-sm text-zinc-500 line-through">
                {formatRUB(product.original_price)}
              </span>
            )}
          </div>
          {exchangeRate > 0 && (
            <span className="text-xs text-zinc-400">{formatUSDT(product.price, exchangeRate)}</span>
          )}
        </div>
      </div>
    </Link>
  );
}
