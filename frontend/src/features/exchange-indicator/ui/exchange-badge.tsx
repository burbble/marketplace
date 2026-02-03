"use client";

import { useState, useEffect } from "react";
import { getExchangeRate } from "@/entities/exchange/api";
import { Badge } from "@/shared/ui/badge";
import { useTranslation } from "@/shared/i18n";

export function ExchangeBadge() {
  const { t } = useTranslation();
  const [rate, setRate] = useState<number | null>(null);

  useEffect(() => {
    let active = true;

    async function fetchRate() {
      try {
        const data = await getExchangeRate();
        if (active) setRate(data.rate);
      } catch {
        if (active) setRate(null);
      }
    }

    fetchRate();
    const interval = setInterval(fetchRate, 60_000);

    return () => {
      active = false;
      clearInterval(interval);
    };
  }, []);

  if (rate === null) {
    return (
      <Badge variant="warning">
        <span className="h-1.5 w-1.5 rounded-full bg-current" />
        {t("exchange.unavailable")}
      </Badge>
    );
  }

  return (
    <Badge variant="success">
      <span className="h-1.5 w-1.5 rounded-full bg-current" />
      {t("exchange.rate", { rate: rate.toFixed(2) })}
    </Badge>
  );
}
