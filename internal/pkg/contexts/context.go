package contexts

import "context"

type ContextKey string

const (
	ContextKeyContext      = "ctx"
	ContextKeyRequestID    = "request_id"
	ContextKeyLogger       = "logger"
	ContextKeyUser         = "user"
	ContextKeyUserID       = "user_id"
	ContextKeyUserName     = "user_name"
	ContextKeyTx           = "tx"
	ContextKeyTelegramID   = "telegram_id"
	ContextKeyTelegramName = "telegram_name"
)

func Set(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, ContextKey(key), value)
}

func Get[T any](ctx context.Context, key string) (T, bool) {
	value := ctx.Value(ContextKey(key))
	if value == nil {
		var t T
		return t, false
	}

	return value.(T), true
}
