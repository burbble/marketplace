"use client";

const SORT_OPTIONS = [
  { label: "Newest", value: "created_at:desc" },
  { label: "Price: Low to High", value: "price:asc" },
  { label: "Price: High to Low", value: "price:desc" },
  { label: "Name: A-Z", value: "name:asc" },
  { label: "Name: Z-A", value: "name:desc" },
  { label: "Brand", value: "brand:asc" },
];

export function SortSelect({
  value,
  onChange,
}: {
  value: string;
  onChange: (v: string) => void;
}) {
  return (
    <select
      value={value || "created_at:desc"}
      onChange={(e) => onChange(e.target.value)}
      className="rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-sm text-zinc-200 outline-none transition-colors focus:border-zinc-500"
    >
      {SORT_OPTIONS.map((opt) => (
        <option key={opt.value} value={opt.value}>
          {opt.label}
        </option>
      ))}
    </select>
  );
}
