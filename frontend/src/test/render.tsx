import { render, type RenderOptions } from "@testing-library/react";
import { I18nProvider } from "@/shared/i18n";

export function renderWithI18n(ui: React.ReactElement, options?: RenderOptions) {
  return render(ui, {
    wrapper: ({ children }) => <I18nProvider>{children}</I18nProvider>,
    ...options,
  });
}
