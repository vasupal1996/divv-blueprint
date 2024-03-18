package logger

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type ApplicationLoggerOpts struct {
	ConsoleWriter io.Writer
	FileWriter    io.Writer
	Config        *ApplicationLoggerConfig
}

type ApplicationLogger struct{}

type ZerlogConfig struct {
	EnableStackTrace bool
	EnableCaller     bool
	Component        string
}

type HookConfig struct {
	EnableHook        bool
	EnableTracingHook bool
	EnableSentryHook  bool
}

type ApplicationLoggerConfig struct {
	ZerlogConfig
	HookConfig
}

func (al *ApplicationLogger) getZerolog(w zerolog.LevelWriter, config *ApplicationLoggerConfig) *zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "msg"

	zlog := zerolog.New(w).With().Timestamp().Logger()

	if config.EnableHook {
		if config.EnableSentryHook {
			zlog = zlog.Hook(&SentryHook{})
		}
		if config.EnableTracingHook {
			zlog = zlog.Hook(&TracingHook{})
		}
	}

	if config.EnableStackTrace {
		zlog = zlog.With().Stack().Logger()
	}

	if config.EnableCaller {
		zlog = zlog.With().Caller().Logger()
	}

	if config.Component != "" {
		zlog = zlog.With().Str("module", config.Component).Logger()
	}

	return &zlog
}

func (al *ApplicationLogger) Setup(opts *ApplicationLoggerOpts) *zerolog.Logger {
	var writers []io.Writer

	// Setting up kafka writer if True.
	// if kw != nil {
	// 	wr := diode.NewWriter(kw, 1000, 10*time.Millisecond, func(missed int) {
	// 		fmt.Printf("Logger Dropped %d messages", missed)
	// 	})
	// 	writers = append(writers, wr)
	// }

	// Setting up console writer if True.
	if opts.ConsoleWriter != nil {
		writers = append(writers, opts.ConsoleWriter)
	}

	// Setting up file writer is True.
	// if fw != nil {
	// 	wr := diode.NewWriter(fw, 1000, 10*time.Millisecond, func(missed int) {
	// 		fmt.Printf("Logger Dropped %d messages", missed)
	// 	})
	// 	writers = append(writers, wr)
	// }

	mw := zerolog.MultiLevelWriter(writers...)
	zlog := al.getZerolog(mw, opts.Config)
	return zlog

}

func (al *ApplicationLogger) CreateSubLogger(logger *zerolog.Logger, name string) *zerolog.Logger {
	l := logger.With().Str("module", name).Logger()
	return &l
}
