package zapx

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestInit_DevMode(t *testing.T) {
	lg, err := Init(Dev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lg == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestInit_ProdMode(t *testing.T) {
	lg, err := Init(Prod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lg == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestInit_NopeMode(t *testing.T) {
	lg, err := Init(Nope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lg == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestInit_UnknownMode(t *testing.T) {
	lg, err := Init("unknown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lg == nil {
		t.Fatal("expected non-nil logger for unknown mode")
	}
}

func TestInit_WithFields(t *testing.T) {
	lg, err := Init(Nope, zap.String("service", "test"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lg == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestWithLogger_AndL(t *testing.T) {
	lg, _ := Init(Nope)
	ctx := WithLogger(context.Background(), lg)

	got := L(ctx)
	if got != lg {
		t.Error("expected logger from context to match")
	}
}

func TestL_NilContext(t *testing.T) {
	got := L(nil) //nolint:staticcheck // intentional nil
	if got == nil {
		t.Error("expected non-nil logger from nil context")
	}
}

func TestL_EmptyContext(t *testing.T) {
	got := L(context.Background())
	if got == nil {
		t.Error("expected non-nil logger from empty context")
	}
}

func TestWithRID_AndGetRID(t *testing.T) {
	ctx := WithRID(context.Background(), "req-123")
	rid := GetRID(ctx)
	if rid != "req-123" {
		t.Errorf("expected 'req-123', got %q", rid)
	}
}

func TestGetRID_NilContext(t *testing.T) {
	rid := GetRID(nil) //nolint:staticcheck // intentional nil
	if rid != "" {
		t.Errorf("expected empty string, got %q", rid)
	}
}

func TestGetRID_NoValue(t *testing.T) {
	rid := GetRID(context.Background())
	if rid != "" {
		t.Errorf("expected empty string, got %q", rid)
	}
}

func TestL_WithRID(t *testing.T) {
	_, _ = Init(Nope)
	ctx := WithRID(context.Background(), "test-rid")
	got := L(ctx)
	if got == nil {
		t.Error("expected non-nil logger with RID context")
	}
}

func TestLogIfErr_WithError(t *testing.T) {
	_, _ = Init(Nope)
	// should not panic
	LogIfErr(context.Background(), nil, "no error")
	LogIfErr(context.Background(), context.Canceled, "canceled")
}

func TestInfo_DoesNotPanic(t *testing.T) {
	_, _ = Init(Nope)
	ctx := context.Background()
	Info(ctx, "test message")
	Warn(ctx, "test warning")
	Error(ctx, "test error")
	Debug(ctx, "test debug")
}
