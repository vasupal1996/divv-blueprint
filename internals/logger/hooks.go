package logger

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type TracingHook struct{}

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	tracingID := getSpanIdFromContext(ctx) // as per your tracing framework
	e.Interface("tracingID", tracingID)
}

func getSpanIdFromContext(ctx context.Context) interface{} {
	return ctx.Value("tracing-id")
}

type SentryHook struct{}

func (h SentryHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.WarnLevel {
		fmt.Println("SENDING TO SENTRY")
		// TODO: Add code to send the log to sentry
	}
}
