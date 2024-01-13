package internals

import (
	"divvy-go-app/internals/logger"
	"divvy-go-app/internals/schema"

	"github.com/rs/zerolog"
)

type App interface {
	Start()
	Close()
}

type AppImpl struct {
	DoneCh chan struct{}
	Logger *zerolog.Logger
}

func (a *AppImpl) Start() {
	a.setupLogger()
}

func (a *AppImpl) Close() {
}

func (a *AppImpl) setupLogger() {
	l := logger.AppLogger{}
	a.Logger = l.Setup(&schema.CreateAppLoggerConfig{ConsoleWriter: logger.NewZeroLogConsoleWriter()})
	a.Logger.Debug().Msg("Logger initiated")
}
