package logrus

import (
	"context"
	"github.com/dougefr/go-clean-arch/interface/iinfra"
	log "github.com/sirupsen/logrus"
	"runtime"
)

type loggerProvider struct {
	l *log.Logger
}

// NewLog ...
func NewLog(logLevel string) (l iinfra.LogProvider, err error) {
	var level log.Level
	level, err = log.ParseLevel(logLevel)
	if err != nil {
		return
	}

	logger := log.New()
	logger.SetLevel(level)
	logger.SetFormatter(&log.JSONFormatter{})
	l = loggerProvider{
		l: logger,
	}

	return
}

func (l loggerProvider) Info(ctx context.Context, message string, attrs ...iinfra.LogAttrs) {
	if l.l.IsLevelEnabled(log.InfoLevel) {
		if a, ok := ctx.Value(iinfra.ContextKeyGlobalLogAttrs).(iinfra.LogAttrs); ok {
			attrs = append(attrs, a)
		}

		attrs = append(attrs, iinfra.LogAttrs{"func": trace()})
		l.l.WithFields(mergeAttrs(attrs)).Info(message)
	}
}

func (l loggerProvider) Error(ctx context.Context, message string, attrs ...iinfra.LogAttrs) {
	if l.l.IsLevelEnabled(log.ErrorLevel) {
		if a, ok := ctx.Value(iinfra.ContextKeyGlobalLogAttrs).(iinfra.LogAttrs); ok {
			attrs = append(attrs, a)
		}

		attrs = append(attrs, iinfra.LogAttrs{"func": trace()})
		l.l.WithFields(mergeAttrs(attrs)).Error(message)
	}
}

func (l loggerProvider) Debug(ctx context.Context, message string, attrs ...iinfra.LogAttrs) {
	if l.l.IsLevelEnabled(log.DebugLevel) {
		if a, ok := ctx.Value(iinfra.ContextKeyGlobalLogAttrs).(iinfra.LogAttrs); ok {
			attrs = append(attrs, a)
		}

		attrs = append(attrs, iinfra.LogAttrs{"func": trace()})
		l.l.WithFields(mergeAttrs(attrs)).Debug(message)
	}
}

func (l loggerProvider) Warn(ctx context.Context, message string, attrs ...iinfra.LogAttrs) {
	if l.l.IsLevelEnabled(log.WarnLevel) {
		if a, ok := ctx.Value(iinfra.ContextKeyGlobalLogAttrs).(iinfra.LogAttrs); ok {
			attrs = append(attrs, a)
		}

		attrs = append(attrs, iinfra.LogAttrs{"func": trace()})
		l.l.WithFields(mergeAttrs(attrs)).Warn(message)
	}
}

func mergeAttrs(attrs []iinfra.LogAttrs) (attr map[string]interface{}) {
	attr = make(iinfra.LogAttrs)

	if attrs == nil {
		return
	}

	for _, a := range attrs {
		for key, value := range a {
			attr[key] = value
		}
	}

	return
}

func trace() string {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	return runtime.FuncForPC(pc[1]).Name()
}
