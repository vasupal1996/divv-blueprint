package router

import (
	"go-app/internals/config"
	"go-app/internals/logger"
	"go-app/schema"
	"go-app/service"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
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
	Logger    *zerolog.Logger
	Config    *config.RouterConfig
	Validator *CustomValidator

	DemoService service.DemoService
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
		Config:      opts.RouterConfig,
		Validator:   NewValidator(),
		DemoService: opts.DemoService,
	}

	r.enableMiddlewares(&middlewareConfig{logger: rr})
	r.RegisterRoutes()
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

// CustomValidator container validator library and transaltor
type CustomValidator struct {
	validator  *validator.Validate
	translator *ut.Translator
}

// NewValidation create new Validator struct instance
func NewValidator() *CustomValidator {
	v := &CustomValidator{
		validator: validator.New(),
	}
	trans := initializeTranslation(v.validator)
	v.translator = trans
	registerFunc(v.validator)
	return v
}

// Initialize initializes and returns the UniversalTranslator instance for the application
func initializeTranslation(validate *validator.Validate) *ut.Translator {

	// initialize translator
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")
	// initialize translations
	en_translations.RegisterDefaultTranslations(validate, trans)
	return &trans
}

func registerFunc(validate *validator.Validate) {
	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate validates the struct
// Note: do not pass slice of struct
func (cv *CustomValidator) Validate(form interface{}) []ErrorResp {
	var validationErrs []ErrorResp
	if err := cv.validator.Struct(form); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var ve ErrorResp
			ve.ErrField = strings.SplitAfterN(e.Namespace(), ".", 2)[1]
			ve.ErrMsg = e.Translate(*cv.translator)
			ve.ErrCode = "ValidationErr"
			validationErrs = append(validationErrs, ve)
		}
	}
	return validationErrs
}
