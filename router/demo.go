package router

import (
	"context"
	"fmt"
	"go-app/schema"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) HelloWorldHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := map[string]string{
		"message": "Hello World",
	}
	r.Logger.Info().Ctx(ctx).Msg("here-1")
	r.DemoService.DemoFunc(ctx)
	return c.JSON(resp)
}

func (r *Router) InternalServerHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	r.Logger.Info().Ctx(ctx).Msg("here-1")
	var x []string
	fmt.Println(x[0])
	return c.Status(http.StatusInternalServerError).JSON(NewErrResponse(false, NewErr("InternalServerErr", "oops something went wrong")))
}

func (r *Router) BadRequestHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	r.Logger.Info().Ctx(ctx).Msg("here-1")
	return c.Status(fiber.StatusBadRequest).JSON(NewErrResponse(false, NewErr("BadRequest", "request failed")))
}

func (r *Router) BadRequestWithSentryWarningHandler(c *fiber.Ctx) error {
	r.Logger.Warn().Ctx(context.WithValue(c.Context(), "meta", "some useful information for debugging")).Msg("new warning message for sentry")
	return c.Status(fiber.StatusBadRequest).JSON(NewErrResponse(false, NewErr("BadRequest", "request failed with sentry")))
}

func (r *Router) BadRequestWithSentryWarningInsideServiceHandler(c *fiber.Ctx) error {
	r.DemoService.SentryDemoFunc(c.Context())
	return fiber.NewError(fiber.StatusBadRequest, "request failed with sentry within service")
}

func (r *Router) InsertOneHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	s := new(schema.InsertOneOpts)
	if err := DecodeJSONBody(c, s); err != nil {
		return c.Status(http.StatusBadRequest).JSON(NewErrResponse(false, err.(ErrorResp)))
	}
	if err := r.Validator.Validate(s); err != nil {
		return c.Status(http.StatusBadRequest).JSON(NewErrResponse(false, err...))
	}
	r.Logger.Info().Ctx(ctx).Msg("here-2")
	id, _ := r.DemoService.InsertOne(ctx, s)
	return c.Status(http.StatusOK).JSON(NewJSONResp(true, fiber.Map{"id": id.Hex()}))
}
