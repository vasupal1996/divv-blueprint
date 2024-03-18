package ws

import (
	"context"
	"divvy-go-app/internals/config"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"

	"github.com/rs/zerolog"
)

type Server interface {
	Start()
	Close() error
}

type FiberServer struct {
	*fiber.App
	Ctx    context.Context
	Worker *sync.WaitGroup
	Logger *zerolog.Logger
	Config *config.WebServerConfig
}

type FiberServerOpts struct {
	FiberApp *fiber.App
	Ctx      context.Context
	Worker   *sync.WaitGroup
	Config   *config.WebServerConfig
	Logger   *zerolog.Logger
}

func NewWebServer(opts *FiberServerOpts) Server {
	s := FiberServer{
		Ctx:    opts.Ctx,
		App:    opts.FiberApp,
		Config: opts.Config,
		Logger: opts.Logger,
		Worker: opts.Worker,
	}
	return &s
}

func (fs *FiberServer) Start() {
	err := fs.Listen(fmt.Sprintf(`%s:%d`, fs.Config.Host, fs.Config.Port))
	if err != nil {
		fs.Logger.Fatal().Err(err).Msg("failed to start server")
	}

}

func (fs *FiberServer) Close() error {
	err := fs.App.Shutdown()
	if err != nil {
		fs.Logger.Err(err).Msg("error while closing webserver")
	} else {
		fs.Logger.Debug().Msg("webserver closed")
	}
	return err
}
