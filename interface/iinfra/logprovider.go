package iinfra

import "context"

// ContextKeyGlobalLogAttrs ...
const ContextKeyGlobalLogAttrs string = "ContextKeyGlobalLogAttrs"

// LogProvider ...
type (
	LogProvider interface {
		Info(ctx context.Context, message string, attrs ...LogAttrs)
		Error(ctx context.Context, message string, attrs ...LogAttrs)
		Debug(ctx context.Context, message string, attrs ...LogAttrs)
		Warn(ctx context.Context, message string, attrs ...LogAttrs)
	}

	// LogAttrs ...
	LogAttrs map[string]interface{}
)
