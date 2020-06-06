package iinfra

import "context"

// ContextKeyTx ...
const ContextKeyGlobalLogAttrs string = "ContextKeyGlobalLogAttrs"

// LogProvider ...
type LogProvider interface {
	Info(ctx context.Context, message string, attrs ...LogAttrs)
	Error(ctx context.Context, message string, attrs ...LogAttrs)
	Debug(ctx context.Context, message string, attrs ...LogAttrs)
	Warn(ctx context.Context, message string, attrs ...LogAttrs)
}

type LogAttrs map[string]interface{}
