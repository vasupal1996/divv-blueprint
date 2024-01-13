package logger

import (
	"divvy-go-app/internals/schema"
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type AppLogger struct{}

func (al *AppLogger) getZerolog(w zerolog.LevelWriter) *zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "msg"

	zlog := zerolog.New(w).With().Timestamp().Stack().Caller().Logger()
	zlog = zlog.Hook(&TracingHook{})
	zlog = zlog.Hook(&SentryHook{})

	return &zlog
}

func (al *AppLogger) Setup(c *schema.CreateAppLoggerConfig) *zerolog.Logger {
	var writers []io.Writer

	// Setting up kafka writer if True.
	// if kw != nil {
	// 	wr := diode.NewWriter(kw, 1000, 10*time.Millisecond, func(missed int) {
	// 		fmt.Printf("Logger Dropped %d messages", missed)
	// 	})
	// 	writers = append(writers, wr)
	// }

	// Setting up console writer if True.
	if c.ConsoleWriter != nil {
		writers = append(writers, c.ConsoleWriter)
	}

	// Setting up file writer is True.
	// if fw != nil {
	// 	wr := diode.NewWriter(fw, 1000, 10*time.Millisecond, func(missed int) {
	// 		fmt.Printf("Logger Dropped %d messages", missed)
	// 	})
	// 	writers = append(writers, wr)
	// }

	mw := zerolog.MultiLevelWriter(writers...)
	zlog := al.getZerolog(mw)
	return zlog

}

func (al *AppLogger) CreateSubLogger(logger *zerolog.Logger, name string) zerolog.Logger {
	return logger.With().Str("module", name).Logger()
}
