import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Spinner } from "./spinner";

describe("Spinner", () => {
  it("renders with status role", () => {
    render(<Spinner />);
    expect(screen.getByRole("status")).toBeInTheDocument();
  });

  it("has accessible loading text", () => {
    render(<Spinner />);
    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("applies custom className", () => {
    render(<Spinner className="mt-4" />);
    const el = screen.getByRole("status");
    expect(el).toHaveClass("mt-4");
  });

  it("has spin animation class", () => {
    render(<Spinner />);
    const el = screen.getByRole("status");
    expect(el).toHaveClass("animate-spin");
  });
});
