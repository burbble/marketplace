package zapx

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey int

const (
	keyRID ctxKey = 100 + iota
	keyLogger
)

const (
	Nope = "nope"
	Prod = "prod"
	Dev  = "dev"
)

func Init(mode string, fields ...zap.Field) (*zap.Logger, error) {
	var (
		lg  *zap.Logger
		err error
	)

	switch mode {
	case Dev:
		lg, err = zap.NewDevelopment()
	case Prod:
		cfg := zap.NewProductionConfig()
		lg, err = cfg.Build(zap.AddStacktrace(zapcore.FatalLevel))
	case Nope:
		lg = zap.NewNop()
	default:
		lg = zap.NewNop()
	}

	if err != nil {
		return nil, err
	}

	if len(fields) > 0 {
		lg = lg.With(fields...)
	}

	_ = zap.ReplaceGlobals(lg)

	return lg, nil
}

func WithLogger(ctx context.Context, lg *zap.Logger) context.Context {
	return context.WithValue(ctx, keyLogger, lg)
}

func WithRID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, keyRID, rid)
}

func GetRID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if v := ctx.Value(keyRID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

func LG(c *gin.Context) *zap.Logger {
	if c != nil && c.Request != nil {
		return L(c.Request.Context())
	}

	return zap.L()
}

func L(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.L()
	}

	if v := ctx.Value(keyLogger); v != nil {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}

	if rid := GetRID(ctx); rid != "" {
		return zap.L().With(zap.String("request_id", rid))
	}

	return zap.L()
}

func Info(ctx context.Context, msg string, fields ...zap.Field)  { L(ctx).Info(msg, fields...) }
func Warn(ctx context.Context, msg string, fields ...zap.Field)  { L(ctx).Warn(msg, fields...) }
func Error(ctx context.Context, msg string, fields ...zap.Field) { L(ctx).Error(msg, fields...) }
func Debug(ctx context.Context, msg string, fields ...zap.Field) { L(ctx).Debug(msg, fields...) }

func LogIfErr(ctx context.Context, err error, msg string, fields ...zap.Field) {
	if err != nil {
		L(ctx).Error(msg, append(fields, zap.Error(err))...)
	}
}
