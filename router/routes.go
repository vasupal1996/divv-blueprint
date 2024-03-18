package router

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func (r *Router) registerRoutes() {
	r.App.Get("/metrics", monitor.New(monitor.Config{Refresh: time.Second * 10}))
	r.App.Get("/", r.HelloWorldHandler)
	r.App.Get("/internal-error", r.InternalServerHandler)
	r.App.Get("/bad-request", r.BadRequestHandler)
	r.App.Get("/bad-request-2", r.BadRequestWithSentryWarningHandler)
	r.App.Get("/bad-request-3", r.BadRequestWithSentryWarningInsideServiceHandler)
	r.App.Get("/insert", r.InsertOneHandler)
}
