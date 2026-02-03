import { cn } from "@/shared/lib/format";

type Variant = "default" | "success" | "warning" | "error";

const variants: Record<Variant, string> = {
  default: "bg-zinc-800 text-zinc-200",
  success: "bg-emerald-900/50 text-emerald-400",
  warning: "bg-amber-900/50 text-amber-400",
  error: "bg-red-900/50 text-red-400",
};

export function Badge({
  children,
  variant = "default",
  className,
}: {
  children: React.ReactNode;
  variant?: Variant;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium",
        variants[variant],
        className
      )}
    >
      {children}
    </span>
  );
}
