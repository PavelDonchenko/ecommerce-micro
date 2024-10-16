package logger

import "context"

type ctxKey int

const (
	traceKey ctxKey = iota + 1
)

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey, traceID)
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(traceKey).(string)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}

	return v
}
