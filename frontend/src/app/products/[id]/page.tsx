"use client";

import { useState, useEffect, use } from "react";
import Link from "next/link";
import { getProductById } from "@/entities/product/api";
import { getExchangeRate } from "@/entities/exchange/api";
import type { Product } from "@/entities/product/model";
import { formatRUB, formatUSDT, resolveImageUrl } from "@/shared/lib/format";
import { Spinner } from "@/shared/ui/spinner";
import { Badge } from "@/shared/ui/badge";
import { useTranslation } from "@/shared/i18n";

export default function ProductPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const { t } = useTranslation();
  const [product, setProduct] = useState<Product | null>(null);
  const [exchangeRate, setExchangeRate] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function load() {
      try {
        const [p, rate] = await Promise.allSettled([getProductById(id), getExchangeRate()]);
        if (p.status === "fulfilled") setProduct(p.value);
        else setError(t("product.notFound"));
        if (rate.status === "fulfilled") setExchangeRate(rate.value.rate);
      } catch {
        setError(t("product.loadError"));
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [id, t]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Spinner className="h-8 w-8" />
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6">
        <div className="flex flex-col items-center gap-4 py-20 text-zinc-500">
          <p className="text-lg">{error || t("product.notFound")}</p>
          <Link
            href="/"
            className="rounded-lg border border-zinc-700 bg-zinc-800 px-4 py-2 text-sm text-zinc-300 transition-colors hover:border-zinc-600"
          >
            {t("product.backToCatalog")}
          </Link>
        </div>
      </div>
    );
  }

  const discount = product.original_price - product.price;

  return (
    <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6">
      <Link
        href="/"
        className="mb-6 inline-flex items-center gap-2 text-sm text-zinc-400 transition-colors hover:text-zinc-200"
      >
        <svg
          className="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
        {t("product.backToCatalog")}
      </Link>

      <div className="grid gap-8 md:grid-cols-2">
        <div className="relative aspect-square overflow-hidden rounded-2xl border border-zinc-800 bg-zinc-900">
          {product.image_url ? (
            <img
              src={resolveImageUrl(product.image_url)}
              alt={product.name}
              className="h-full w-full object-cover"
            />
          ) : (
            <div className="flex h-full items-center justify-center text-zinc-600">
              {t("product.noImage")}
            </div>
          )}
        </div>

        <div className="flex flex-col gap-6">
          <div className="flex flex-col gap-2">
            <div className="flex flex-wrap items-center gap-2">
              <Badge>{product.brand}</Badge>
              <Badge variant="default">
                {t("product.sku")}: {product.sku}
              </Badge>
            </div>
            <h1 className="text-2xl font-bold text-white sm:text-3xl">{product.name}</h1>
          </div>

          <div className="flex flex-col gap-3 rounded-xl border border-zinc-800 bg-zinc-900 p-6">
            <div className="flex items-baseline gap-3">
              <span className="text-3xl font-bold text-white">{formatRUB(product.price)}</span>
              {discount > 0 && (
                <span className="text-lg text-zinc-500 line-through">
                  {formatRUB(product.original_price)}
                </span>
              )}
            </div>

            {discount > 0 && <Badge variant="success">-{formatRUB(discount)}</Badge>}

            {exchangeRate > 0 && (
              <p className="text-sm text-zinc-400">≈ {formatUSDT(product.price, exchangeRate)}</p>
            )}
          </div>

          {product.description && (
            <div className="flex flex-col gap-3 rounded-xl border border-zinc-800 bg-zinc-900 p-6">
              <h2 className="text-sm font-medium text-zinc-400">{t("product.description")}</h2>
              <p className="text-sm leading-relaxed text-zinc-300 whitespace-pre-line">
                {product.description}
              </p>
            </div>
          )}

          <div className="flex flex-col gap-3 rounded-xl border border-zinc-800 bg-zinc-900 p-6">
            <h2 className="text-sm font-medium text-zinc-400">{t("product.details")}</h2>
            <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm">
              <dt className="text-zinc-500">{t("product.brand")}</dt>
              <dd className="text-zinc-200">{product.brand}</dd>
              <dt className="text-zinc-500">{t("product.sku")}</dt>
              <dd className="text-zinc-200">{product.sku}</dd>
              <dt className="text-zinc-500">{t("product.externalId")}</dt>
              <dd className="text-zinc-200">{product.external_id}</dd>
            </dl>
          </div>
        </div>
      </div>
    </div>
  );
}
