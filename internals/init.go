package internals

import (
	"context"
	"divvy-go-app/internals/config"
	"divvy-go-app/internals/db"
	"divvy-go-app/internals/logger"
	"divvy-go-app/internals/mongodb"

	"divvy-go-app/internals/ws"
	"divvy-go-app/router"
	"divvy-go-app/service"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber"
	"github.com/rs/zerolog"
)

type App interface {
	Start()
	Close()
}

type AppImpl struct {
	AbstractLogger *logger.ApplicationLogger
	Config         *config.Config
	Ctx            context.Context
	Logger         *zerolog.Logger
	DB             db.DB
	Service        service.Service
	WebServer      ws.Server
}

func CreateNewApp(ctx context.Context) App {
	a := AppImpl{
		Ctx: ctx,
	}
	return &a
}

func (a *AppImpl) Start() {
	a.setupLogger()
	a.getConfig()
	a.setupDB()
	a.setupService()
	a.setupWebServer()
	a.setupSentry()
}

func (a *AppImpl) Close() {
	<-a.Ctx.Done()

	// Closing down all the components
	a.WebServer.Close()
	a.Service.Close()
	a.DB.MongoDB().Close()
	a.Logger.Debug().Msg("app gracefully closed")
}

func (a *AppImpl) setupLogger() {
	a.AbstractLogger = &logger.ApplicationLogger{}
	a.Logger = a.AbstractLogger.Setup(&logger.ApplicationLoggerOpts{
		ConsoleWriter: logger.NewZeroLogConsoleWriter(),
		Config: &logger.ApplicationLoggerConfig{
			ZerlogConfig: logger.ZerlogConfig{
				EnableStackTrace: true,
				EnableCaller:     true,
			},
			HookConfig: logger.HookConfig{
				EnableHook:        true,
				EnableTracingHook: true,
			},
		},
	})
}

func (a *AppImpl) getConfig() {
	c := config.GetConfigFromFile()
	config.WatchConfigChanges(a.Logger, c)
	a.Config = c
}

func (a *AppImpl) setupDB() {
	a.DB = db.NewDB(&db.DBOpts{
		MongoDB: a.setupMongoDB(),
	})
}

func (a *AppImpl) setupMongoDB() mongodb.MongoDB {
	mongodb, err := mongodb.NewMongoDB(&mongodb.MongoDBOpts{
		Config: a.Config.MongoDBConfig,
		Logger: a.AbstractLogger.CreateSubLogger(a.Logger, "mongodb"),
		Ctx:    a.Ctx,
	})

	if err != nil {
		a.Logger.Fatal().Err(err).Msg("failed to setup mongodb")
	}
	return mongodb
}

func (a *AppImpl) setupRouter() *router.Router {
	router := router.NewRouter(&router.RouterOpts{
		AbstractLogger: a.AbstractLogger,
		DemoService:    a.Service.GetDemoService(),
		RouterConfig:   a.Config.RouterConfig,
	})
	return router
}

func (a *AppImpl) setupWebServer() {
	router := a.setupRouter()
	a.WebServer = ws.NewWebServer(&ws.FiberServerOpts{
		FiberApp: router.App,
		Ctx:      a.Ctx,
		Config:   a.Config.WebServerConfig,
		Logger:   a.AbstractLogger.CreateSubLogger(a.Logger, "ws"),
	})

	go a.WebServer.Start()
}

func (a *AppImpl) setupService() {
	a.Service = service.NewService(&service.ServiceOpts{
		Ctx:            a.Ctx,
		Logger:         a.Logger,
		AbstractLogger: a.AbstractLogger,
		Config:         a.Config.AppConfig.ServiceConfig,
		DB:             a.DB,
	})
}

func (a *AppImpl) setupSentry() {
	if a.Config.SentryConfig.EnableSentry {
		_ = sentry.Init(sentry.ClientOptions{
			Dsn: a.Config.SentryConfig.Host,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				if hint.Context != nil {
					if _, ok := hint.Context.Value(sentry.RequestContextKey).(*fiber.Ctx); ok {
						// You have access to the original Context if it panicked
						fmt.Println(ok)
					}
				}
				return event
			},
			EnableTracing:    true,
			AttachStacktrace: true,
		})
	}
}
