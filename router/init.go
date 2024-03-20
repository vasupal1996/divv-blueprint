package router

import (
	"go-app/internals/config"
	"go-app/internals/logger"
	"go-app/schema"
	"go-app/service"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
)

var RequestFieldsToLog = []string{
	"time",
	"referer",
	"protocol",
	"pid",
	"port",
	"ip",
	"ips",
	"host",
	"path",
	"url",
	"ua",
	"latency",
	"status",
	"resBody",
	"queryParams",
	"body",
	"bytesReceived",
	"bytesSent",
	"route",
	"method",
	"requestId",
	"error",
	"reqHeaders",
}

type Router struct {
	*fiber.App
	Logger *zerolog.Logger
	Config *config.RouterConfig

	demoService service.DemoService
}

type RouterOpts struct {
	AbstractLogger *logger.ApplicationLogger
	DemoService    service.DemoService
	RouterConfig   *config.RouterConfig
}

type middlewareConfig struct {
	logger *zerolog.Logger
}

func NewRouter(opts *RouterOpts) *Router {
	lr := opts.AbstractLogger.Setup(&logger.ApplicationLoggerOpts{
		ConsoleWriter: logger.NewZeroLogConsoleWriter(),
		Config: &logger.ApplicationLoggerConfig{
			ZerlogConfig: logger.ZerlogConfig{
				Component:        "router",
				EnableStackTrace: true,
				EnableCaller:     true,
			},
			HookConfig: logger.HookConfig{
				EnableHook:        true,
				EnableTracingHook: true,
				EnableSentryHook:  true,
			},
		},
	})

	rr := opts.AbstractLogger.Setup(&logger.ApplicationLoggerOpts{
		ConsoleWriter: logger.NewZeroLogConsoleWriter(),
		Config: &logger.ApplicationLoggerConfig{
			ZerlogConfig: logger.ZerlogConfig{
				Component: "requests",
			},
			HookConfig: logger.HookConfig{
				EnableHook:        true,
				EnableTracingHook: true,
			},
		},
	})

	fiberConfig := fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	}

	r := Router{
		App:         fiber.New(fiberConfig),
		Logger:      lr,
		demoService: opts.DemoService,
		Config:      opts.RouterConfig,
	}

	r.enableMiddlewares(&middlewareConfig{logger: rr})
	r.registerRoutes()
	return &r
}

func (r *Router) enableMiddlewares(config *middlewareConfig) {

	r.App.Use(requestid.New(requestid.Config{
		Header:     "X-Request-ID",
		ContextKey: schema.RequestIDKey,
	}))

	r.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	r.App.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/dont_compress"
		},
	}))

	r.App.Use(helmet.New())

	r.App.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: config.logger,
		Fields: RequestFieldsToLog,
	}))

	if r.Config.EnableSentry {
		r.App.Use(fibersentry.New(fibersentry.Config{
			Repanic:         true,
			WaitForDelivery: true,
		}))
	}

}

// var enhanceSentryEvent = func(c *fiber.Ctx) error {
// 	if hub := fibersentry.GetHubFromContext(c); hub != nil {
// 	}
// 	return c.Next()
// }
