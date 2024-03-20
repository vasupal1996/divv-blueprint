//go:generate $GOPATH/bin/mockgen -destination=../../mock/mock_mongodb.go -package=mock go-app/internals/mongodb MongoDB,Client,Database,Collection,Cursor

package mongodb

import (
	"context"
	"go-app/internals/config"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB interface {
	Close() error
	Cli() Client
}

type MongoDBImpl struct {
	Ctx    context.Context
	Worker *sync.WaitGroup
	Logger *zerolog.Logger
	Config *config.MongoDBConfig
	Client Client
}

type MongoDBOpts struct {
	Ctx    context.Context
	Logger *zerolog.Logger
	Config *config.MongoDBConfig
	Worker *sync.WaitGroup
}

func (mdbi *MongoDBImpl) Cli() Client {
	return mdbi.Client
}

func (mdbi *MongoDBImpl) Close() error {
	err := mdbi.Client.Disconnect(context.TODO())
	if err != nil {
		mdbi.Logger.Err(err).Msg("error while closing mongodb connection")
	} else {
		mdbi.Logger.Debug().Msg("mongodb connection closed")
	}
	// time.Sleep(1 * time.Second)
	return err
}

func getMongoDBClientOpts(c *config.MongoDBConfig) *options.ClientOptions {
	clientOpts := options.Client().ApplyURI(c.ConnectionURL())
	if c.ReadPref != "" {
		mode, err := readpref.ModeFromString(c.ReadPref)
		if err != nil {
			rp, err := readpref.New(mode)
			if err == nil {
				clientOpts = clientOpts.SetReadPreference(rp)
			}
		}
	}
	return clientOpts
}

func NewMongoDB(opts *MongoDBOpts) (MongoDB, error) {
	client, err := NewClient(opts)
	if err != nil {
		return nil, errors.Wrap(err, "connect failed")
	}
	// checking if client is pining or not otherwise return error
	if err := client.Ping(context.TODO()); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	mongodb := MongoDBImpl{Client: client, Ctx: opts.Ctx, Worker: opts.Worker, Logger: opts.Logger, Config: opts.Config}
	return &mongodb, nil
}

func NewMockMongoDB(url string) (MongoDB, error) {
	client, err := NewTestClient(url)
	if err != nil {
		return nil, errors.Wrap(err, "connect failed")
	}
	// checking if client is pining or not otherwise return error
	if err := client.Ping(context.TODO()); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	mongodb := MongoDBImpl{Client: client}
	return &mongodb, nil
}
