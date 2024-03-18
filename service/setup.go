package service

import (
	"context"
	"divvy-go-app/internals/config"

	"github.com/rs/zerolog"
)

type DemoServiceImpl struct {
	Ctx     context.Context
	Logger  *zerolog.Logger
	Config  *config.DemoServiceConfig
	Service Service
}

type DemoServiceOpts struct {
	Ctx     context.Context
	Logger  *zerolog.Logger
	Config  *config.DemoServiceConfig
	Service Service
}

func NewDemoService(opts *DemoServiceOpts) DemoService {
	// l := opts.ServiceConfig.AbstractLogger.CreateSubLogger(opts.ServiceConfig.Logger, "demo-service")
	ds := DemoServiceImpl{
		Ctx:     opts.Ctx,
		Logger:  opts.Logger,
		Config:  opts.Config,
		Service: opts.Service,
	}
	return &ds
}
