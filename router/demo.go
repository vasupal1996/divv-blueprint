package router

import (
	"context"
	"divvy-go-app/schema"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) HelloWorldHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := map[string]string{
		"message": "Hello World",
	}

	r.Logger.Info().Ctx(ctx).Msg("here-1")
	r.demoService.DemoFunc(ctx)

	return c.JSON(resp)
}

func (r *Router) InternalServerHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	r.Logger.Info().Ctx(ctx).Msg("here-1")
	var x []string
	fmt.Println(x[0])
	return c.JSON("If you're seeing this then it's not working")
}

func (r *Router) BadRequestHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	r.Logger.Info().Ctx(ctx).Msg("here-1")
	return fiber.NewError(fiber.StatusBadRequest, "request failed")
}

func (r *Router) BadRequestWithSentryWarningHandler(c *fiber.Ctx) error {
	r.Logger.Warn().Ctx(context.WithValue(c.Context(), "meta", "some useful information for debugging")).Msg("new warning message for sentry")
	return fiber.NewError(fiber.StatusBadRequest, "request failed with sentry")
}

func (r *Router) BadRequestWithSentryWarningInsideServiceHandler(c *fiber.Ctx) error {
	r.demoService.SentryDemoFunc(c.Context())
	return fiber.NewError(fiber.StatusBadRequest, "request failed with sentry within service")
}

func (r *Router) InsertOneHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	s := schema.InsertOneOpts{
		Name: "Vasu",
	}
	r.Logger.Info().Ctx(ctx).Msg("here-2")
	r.demoService.InsertOne(ctx, &s)
	return c.JSON(s)
}
