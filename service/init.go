//go:generate mockgen -destination=../mock/mock_service.go -package=mock go-app/service Service
package service

import (
	"context"
	"go-app/internals/config"
	"go-app/internals/db"
	"go-app/internals/logger"
	"sync"

	"github.com/rs/zerolog"
)

type Service interface {
	Close() error
	GetDemoService() DemoService

	db.DB
}

type ServiceImpl struct {
	AbstractLogger *logger.ApplicationLogger
	Ctx            context.Context
	Config         *config.ServiceConfig
	Logger         *zerolog.Logger
	Sync           *sync.WaitGroup

	db.DB
	DemoService DemoService
}

type ServiceOpts struct {
	AbstractLogger *logger.ApplicationLogger
	Config         *config.ServiceConfig
	Ctx            context.Context
	Logger         *zerolog.Logger
	DB             db.DB
	Sync           *sync.WaitGroup
}

func NewService(opts *ServiceOpts) Service {
	sl := opts.AbstractLogger.Setup(&logger.ApplicationLoggerOpts{
		ConsoleWriter: logger.NewZeroLogConsoleWriter(),
		Config: &logger.ApplicationLoggerConfig{
			ZerlogConfig: logger.ZerlogConfig{
				EnableStackTrace: true,
				EnableCaller:     true,
				Component:        "service",
			},
			HookConfig: logger.HookConfig{
				EnableHook:        true,
				EnableTracingHook: true,
				EnableSentryHook:  true,
			},
		},
	})
	s := ServiceImpl{
		AbstractLogger: opts.AbstractLogger,
		Ctx:            opts.Ctx,
		Logger:         sl,
		Config:         opts.Config,
		DB:             opts.DB,
		Sync:           opts.Sync,
	}
	s.setup(opts)
	return &s
}

func (si *ServiceImpl) Close() error {
	si.Logger.Debug().Msg("services closed")
	return nil
}

func (si *ServiceImpl) GetDemoService() DemoService {
	return si.DemoService
}

func (si *ServiceImpl) setup(opts *ServiceOpts) {
	si.DemoService = NewDemoService(&DemoServiceOpts{
		Config:  opts.Config.DemoServiceConfig,
		Logger:  si.AbstractLogger.CreateSubLogger(si.Logger, "demo-service"),
		Service: si,
	})
}
