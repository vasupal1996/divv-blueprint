package service_test

import (
	"context"
	"go-app/internals/config"
	"go-app/internals/db"
	"go-app/internals/logger"
	"go-app/internals/mongodb"
	"go-app/mock"
	"go-app/service"

	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	memongo "github.com/tryvium-travels/memongo"
	"github.com/tryvium-travels/memongo/memongolog"
)

type TestService struct {
	Service       service.Service
	Ctrl          *gomock.Controller
	MongoDBServer *memongo.Server
}

func (ts *TestService) Clean() {
	ts.Ctrl.Finish()
	ts.MongoDBServer.Stop()
}

func NewTestService(t *testing.T) *TestService {
	ctrl := gomock.NewController(t)

	config := config.GetTestConfigFromFile()
	s := service.ServiceImpl{
		AbstractLogger: &logger.ApplicationLogger{},
		Ctx:            context.TODO(),
		Config:         config.AppConfig.ServiceConfig,
		Logger:         &zerolog.Logger{},
		Sync:           &sync.WaitGroup{},
	}

	s.DemoService = mock.NewMockDemoService(ctrl)

	mongoDBServer := NewMockMongoDB()
	mockMongoDB, err := mongodb.NewMockMongoDB(mongoDBServer.URI())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.DB = db.NewDB(&db.DBOpts{MongoDB: mockMongoDB})
	return &TestService{
		Ctrl:          ctrl,
		Service:       &s,
		MongoDBServer: mongoDBServer,
	}
}

func NewMockMongoDB() *memongo.Server {
	opts := memongo.Options{
		ShouldUseReplica: true,
		MongoVersion:     "6.0.12",
		LogLevel:         memongolog.LogLevelWarn,
	}
	mongoServer, err := memongo.StartWithOptions(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return mongoServer
}
