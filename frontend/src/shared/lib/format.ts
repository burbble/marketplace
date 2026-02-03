export function formatRUB(price: number): string {
  return new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    maximumFractionDigits: 0,
  }).format(price);
}

export function formatUSDT(price: number, rate: number): string {
  if (rate <= 0) return "—";
  const usdt = price / rate;
  return `${usdt.toFixed(2)} USDT`;
}

export function cn(...classes: (string | boolean | undefined | null)[]): string {
  return classes.filter(Boolean).join(" ");
}

export function resolveImageUrl(url: string): string {
  if (!url) return "";
  if (url.startsWith("http")) return url;
  return `https://store77.net${url}`;
}
