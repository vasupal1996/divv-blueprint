package logger

import (
	"context"
	"go-app/schema"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

type TracingHook struct{}

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	tracingID := getSpanIdFromContext(ctx) // as per your tracing framework
	if tracingID != nil {
		e.Interface(schema.RequestIDKey, tracingID)
	}
}

func getSpanIdFromContext(ctx context.Context) interface{} {
	id := ctx.Value(ctx.Value(schema.RequestIDKey))
	return id
}

type SentryHook struct{}

func (h SentryHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	if level == zerolog.WarnLevel {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetContext("ctx", sentry.Context{
				schema.RequestIDKey:   ctx.Value(schema.RequestIDKey),
				schema.SentryExtraCtx: ctx.Value(schema.SentryExtraCtx),
			})
			scope.SetLevel(sentry.LevelWarning)
			sentry.CaptureMessage(msg)
		})
	}
}
