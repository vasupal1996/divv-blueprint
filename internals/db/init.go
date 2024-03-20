package db

import (
	"go-app/internals/mongodb"
)

type DB interface {
	MongoDB() mongodb.MongoDB
}

type DBImpl struct {
	DB mongodb.MongoDB
}

type DBOpts struct {
	MongoDB mongodb.MongoDB
}

func (dbi *DBImpl) MongoDB() mongodb.MongoDB {
	return dbi.DB
}

func NewDB(opts *DBOpts) DB {
	db := DBImpl{
		DB: opts.MongoDB,
	}
	return &db
}

// func mongoDB(ctrl *gomock.Controller) mongodb.MongoDB {
// 	m := mock.NewMockMongoDB(ctrl)
// 	return m
// }

// func NewMock(ctrl *gomock.Controller) DB {
// 	db := DBImpl{
// 		DB: mongoDB(ctrl),
// 	}
// 	return &db
// }
