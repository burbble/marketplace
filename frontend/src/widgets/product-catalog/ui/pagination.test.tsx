import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Pagination } from "./pagination";

describe("Pagination", () => {
  it("returns null when total pages is 1", () => {
    const { container } = render(
      <Pagination page={1} total={10} pageSize={24} onChange={() => {}} />
    );
    expect(container.innerHTML).toBe("");
  });

  it("returns null when total is 0", () => {
    const { container } = render(
      <Pagination page={1} total={0} pageSize={24} onChange={() => {}} />
    );
    expect(container.innerHTML).toBe("");
  });

  it("renders page buttons for small total", () => {
    render(<Pagination page={1} total={120} pageSize={24} onChange={() => {}} />);
    expect(screen.getByText("1")).toBeInTheDocument();
    expect(screen.getByText("5")).toBeInTheDocument();
  });

  it("highlights current page", () => {
    render(<Pagination page={3} total={120} pageSize={24} onChange={() => {}} />);
    const btn = screen.getByText("3");
    expect(btn).toHaveClass("font-medium");
  });

  it("disables prev button on first page", () => {
    render(<Pagination page={1} total={100} pageSize={24} onChange={() => {}} />);
    const buttons = screen.getAllByRole("button");
    expect(buttons[0]).toBeDisabled();
  });

  it("disables next button on last page", () => {
    render(<Pagination page={5} total={100} pageSize={24} onChange={() => {}} />);
    const buttons = screen.getAllByRole("button");
    expect(buttons[buttons.length - 1]).toBeDisabled();
  });

  it("calls onChange with previous page", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<Pagination page={3} total={100} pageSize={24} onChange={onChange} />);
    const buttons = screen.getAllByRole("button");
    await user.click(buttons[0]);
    expect(onChange).toHaveBeenCalledWith(2);
  });

  it("calls onChange with next page", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<Pagination page={3} total={100} pageSize={24} onChange={onChange} />);
    const buttons = screen.getAllByRole("button");
    await user.click(buttons[buttons.length - 1]);
    expect(onChange).toHaveBeenCalledWith(4);
  });

  it("calls onChange with clicked page number", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<Pagination page={1} total={100} pageSize={24} onChange={onChange} />);
    await user.click(screen.getByText("3"));
    expect(onChange).toHaveBeenCalledWith(3);
  });

  it("shows ellipsis for large page counts", () => {
    render(<Pagination page={5} total={240} pageSize={24} onChange={() => {}} />);
    const dots = screen.getAllByText("...");
    expect(dots.length).toBeGreaterThanOrEqual(1);
  });
});
