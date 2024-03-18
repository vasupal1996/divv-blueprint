package mongodb

import (
	"context"

	"reflect"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client interface {
	Database(string, ...*options.DatabaseOptions) Database
	Disconnect(context.Context) error
	StartSession(...*options.SessionOptions) (mongo.Session, error)
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	Ping(context.Context) error
}

type Database interface {
	Collection(string, ...*options.CollectionOptions) Collection
	Client() Client
}

type Cursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
	All(context.Context, interface{}) error
}

type Collection interface {
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (interface{}, error)
	InsertMany(context.Context, []interface{}, ...*options.InsertManyOptions) ([]interface{}, error)
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (int64, error)
	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
	CountDocuments(context.Context, interface{}, ...*options.CountOptions) (int64, error)
	Aggregate(context.Context, interface{}, ...*options.AggregateOptions) (*mongo.Cursor, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type SingleResult interface {
	Decode(interface{}) error
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoCursor struct {
	mc *mongo.Cursor
}

type mongoSession struct {
	mongo.Session
}

type nullawareDecoder struct {
	defDecoder bsoncodec.ValueDecoder
	zeroValue  reflect.Value
}

func (d *nullawareDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() != bson.TypeNull {
		return d.defDecoder.DecodeValue(dctx, vr, val)
	}

	if !val.CanSet() {
		return errors.New("value not settable")
	}
	if err := vr.ReadNull(); err != nil {
		return err
	}
	// Set the zero value of val's type:
	val.Set(d.zeroValue)
	return nil
}

func NewClient(opts *MongoDBOpts) (Client, error) {
	c, err := mongo.Connect(context.TODO(), getMongoDBClientOpts(opts.Config))
	if err != nil {
		return nil, errors.Wrap(err, "connect failed")
	}
	return &mongoClient{cl: c}, err

}

func NewTestClient(url string) (Client, error) {
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		return nil, errors.Wrap(err, "connect failed")
	}
	return &mongoClient{cl: c}, err
}

func (mc *mongoClient) Ping(ctx context.Context) error {
	return mc.cl.Ping(ctx, readpref.Primary())
}

func (mc *mongoClient) Database(dbName string, opts ...*options.DatabaseOptions) Database {
	db := mc.cl.Database(dbName, opts...)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return mc.cl.UseSession(ctx, fn)
}

func (mc *mongoClient) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	session, err := mc.cl.StartSession(opts...)
	return &mongoSession{session}, err
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.cl.Disconnect(ctx)
}

func (md *mongoDatabase) Collection(colName string, opts ...*options.CollectionOptions) Collection {
	collection := md.db.Collection(colName, opts...)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() Client {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return mc.coll.FindOne(ctx, filter)
}

func (mc *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mc.coll.UpdateOne(ctx, filter, update, opts...)
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return id.InsertedID, err
}

func (mc *mongoCollection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) ([]interface{}, error) {
	res, err := mc.coll.InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	count, err := mc.coll.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return mc.coll.Find(ctx, filter, opts...)
}

func (mc *mongoCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return mc.coll.Aggregate(ctx, pipeline, opts...)
}

func (mc *mongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mc.coll.UpdateMany(ctx, filter, update, opts...)
}

func (mc *mongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return mc.coll.CountDocuments(ctx, filter, opts...)
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (mr *mongoCursor) Close(ctx context.Context) error {
	return mr.mc.Close(ctx)
}

func (mr *mongoCursor) Next(ctx context.Context) bool {
	return mr.mc.Next(ctx)
}

func (mr *mongoCursor) Decode(v interface{}) error {
	return mr.mc.Decode(v)
}

func (mr *mongoCursor) All(ctx context.Context, result interface{}) error {
	return mr.mc.All(ctx, result)
}
