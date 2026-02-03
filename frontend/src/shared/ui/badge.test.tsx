import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Badge } from "./badge";

describe("Badge", () => {
  it("renders children", () => {
    render(<Badge>Test Label</Badge>);
    expect(screen.getByText("Test Label")).toBeInTheDocument();
  });

  it("applies default variant classes", () => {
    render(<Badge>Default</Badge>);
    const el = screen.getByText("Default");
    expect(el).toHaveClass("bg-zinc-800");
  });

  it("applies success variant classes", () => {
    render(<Badge variant="success">OK</Badge>);
    const el = screen.getByText("OK");
    expect(el).toHaveClass("text-emerald-400");
  });

  it("applies error variant classes", () => {
    render(<Badge variant="error">Fail</Badge>);
    const el = screen.getByText("Fail");
    expect(el).toHaveClass("text-red-400");
  });

  it("applies warning variant classes", () => {
    render(<Badge variant="warning">Warn</Badge>);
    const el = screen.getByText("Warn");
    expect(el).toHaveClass("text-amber-400");
  });

  it("appends custom className", () => {
    render(<Badge className="my-class">Custom</Badge>);
    const el = screen.getByText("Custom");
    expect(el).toHaveClass("my-class");
  });
});
