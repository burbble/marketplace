import type { Metadata } from "next";
import { NuqsAdapter } from "nuqs/adapters/next/app";
import { Header } from "@/widgets/header/ui/header";
import { I18nProvider } from "@/shared/i18n";
import "./globals.css";

export const metadata: Metadata = {
  title: "Marketplace",
  description: "Product catalog",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="dark">
      <body className="min-h-screen bg-zinc-950 font-sans text-zinc-100 antialiased">
        <I18nProvider>
          <NuqsAdapter>
            <Header />
            <main>{children}</main>
          </NuqsAdapter>
        </I18nProvider>
      </body>
    </html>
  );
}
