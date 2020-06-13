// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package infra

import (
	"context"
	"runtime"

	"github.com/dougefr/go-clean-arch/interface/iinfra"
	log "github.com/sirupsen/logrus"
)

type loggerProvider struct {
	l *log.Logger
}

// NewLogrus ...
func NewLogrus(logLevel string) (l iinfra.LogProvider, err error) {
	var level log.Level
	level, err = log.ParseLevel(logLevel)
	if err != nil {
		return
	}

	logger := log.New()
	logger.SetLevel(level)
	logger.SetFormatter(&log.JSONFormatter{}) // set JSON formatter as default formatter
	l = loggerProvider{
		l: logger,
	}

	return
}

func (l loggerProvider) Info(ctx context.Context, message string, attrs ...iinfra.LogAttrs) {
	if l.l.IsLevelEnabled(log.InfoLevel) {
		go func() {
			attrs = appendGlobalAttrs(ctx, attrs)
			attrs = append(attrs, iinfra.LogAttrs{"func": trace()})
			l.l.WithFields(mergeAttrs(attrs)).Info(message)
		}()
	}
}

func (l loggerProvider) Error(ctx context.Context, message string, attrs ...iinfra.LogAttrs) {
	if l.l.IsLevelEnabled(log.ErrorLevel) {
		attrs = appendGlobalAttrs(ctx, attrs)
		attrs = append(attrs, iinfra.LogAttrs{"func": trace()})
		l.l.WithFields(mergeAttrs(attrs)).Error(message)
	}
}

func (l loggerProvider) Debug(ctx context.Context, message string, attrs ...iinfra.LogAttrs) {
	if l.l.IsLevelEnabled(log.DebugLevel) {
		attrs = appendGlobalAttrs(ctx, attrs)
		attrs = append(attrs, iinfra.LogAttrs{"func": trace()})
		l.l.WithFields(mergeAttrs(attrs)).Debug(message)
	}
}

func (l loggerProvider) Warn(ctx context.Context, message string, attrs ...iinfra.LogAttrs) {
	if l.l.IsLevelEnabled(log.WarnLevel) {
		attrs = appendGlobalAttrs(ctx, attrs)
		attrs = append(attrs, iinfra.LogAttrs{"func": trace()})
		l.l.WithFields(mergeAttrs(attrs)).Warn(message)
	}
}

// appendGlobalAttrs get global attrs from the context
func appendGlobalAttrs(ctx context.Context, attrs []iinfra.LogAttrs) []iinfra.LogAttrs {
	if a, ok := ctx.Value(iinfra.ContextKeyGlobalLogAttrs).(iinfra.LogAttrs); ok {
		attrs = append(attrs, a)
	}
	return attrs
}

// mergeAttrs merge all attrs in just one map
func mergeAttrs(attrs []iinfra.LogAttrs) (attr map[string]interface{}) {
	attr = make(iinfra.LogAttrs)

	if attrs == nil {
		return // empty attrs
	}

	for _, a := range attrs {
		for key, value := range a {
			attr[key] = value
		}
	}

	return
}

// get the name of the function to be logged
func trace() string {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	return runtime.FuncForPC(pc[1]).Name()
}
