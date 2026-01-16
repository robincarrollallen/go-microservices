package trace

import "context"

type traceIDKeyType struct{}

var traceIDKey = traceIDKeyType{}

func TraceIDKey() interface{} {
	return traceIDKey
}

func FromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(traceIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}
