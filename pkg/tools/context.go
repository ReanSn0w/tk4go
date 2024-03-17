package tools

import "context"

var (
	cancelCtxKey = &cancelKey{}
)

type (
	cancelKey struct{}
)

// GlobalCancel creates a new context with a global cancel function
func GlobalCancel(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancelCause(ctx)
	ctx = context.WithValue(ctx, cancelCtxKey, cancel)
	return ctx
}

// GetGlobalCancel returns the global cancel function from the context
func GetGlobalCancel(ctx context.Context) context.CancelCauseFunc {
	cancel, _ := ctx.Value(cancelCtxKey).(context.CancelCauseFunc)
	return cancel
}
