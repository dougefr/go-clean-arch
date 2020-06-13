// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

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
