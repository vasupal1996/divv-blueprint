package router_test

import (
	"go-app/internals/config"
	"go-app/mock"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

type TestRouter struct {
	*fiber.App
	Logger      *zerolog.Logger
	Config      *config.RouterConfig
	Ctrl        *gomock.Controller
	demoService *mock.MockDemoService
}

func (ts *TestRouter) Clean() {
	ts.Ctrl.Finish()
}

func NewRouterTest(t *testing.T) *TestRouter {
	ctrl := gomock.NewController(t)
	r := TestRouter{
		App:    fiber.New(fiber.Config{}),
		Logger: &zerolog.Logger{},
		Config: config.GetTestConfigFromFile().RouterConfig,
		Ctrl:   gomock.NewController(t),
	}
	r.demoService = mock.NewMockDemoService(ctrl)
	return &r
}
