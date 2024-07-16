package constant

type ContextKey string

const (
	ContextKeyContext = "ctx"
	ContextKeyRequestID = "request_id"
	ContextKeyLogger    = "logger"
	ContextKeyUserID    = "user_id"
	ContextKeyUserName    = "user_name"
)
