"use client";

export default function Error({ error, reset }: { error: Error; reset: () => void }) {
  return (
    <div className="flex min-h-[50vh] flex-col items-center justify-center gap-4 px-4">
      <div className="rounded-xl border border-zinc-800 bg-zinc-900 p-8 text-center">
        <h2 className="mb-2 text-lg font-medium text-white">Something went wrong</h2>
        <p className="mb-6 text-sm text-zinc-400">{error.message || "An unexpected error occurred"}</p>
        <button
          onClick={reset}
          className="rounded-lg border border-zinc-700 bg-zinc-800 px-4 py-2 text-sm text-zinc-200 transition-colors hover:border-zinc-600"
        >
          Try again
        </button>
      </div>
    </div>
  );
}
