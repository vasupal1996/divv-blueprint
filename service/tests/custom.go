package service_test

import (
	"context"
	"go-app/internals/mongodb"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Get_DocByFilter(coll mongodb.Collection, filter bson.M, result interface{}) error {
	res := coll.FindOne(context.TODO(), filter)
	err := res.Decode(result)
	return err
}

func Get_DocsByFilter(coll mongodb.Collection, filter bson.M, result interface{}) error {
	cur, _ := coll.Find(context.TODO(), filter)
	err := cur.All(context.TODO(), result)
	return err
}

func Assert_DocCount(t *testing.T, coll mongodb.Collection, filter bson.M, expectedCount int) {
	count, err := coll.CountDocuments(context.TODO(), filter)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedCount, count)
}

func Assert_TimestampDuration(t *testing.T, timestamp1, timestamp2 time.Time, deltaMs ...*time.Duration) {
	if deltaMs != nil {
		assert.WithinDuration(t, timestamp1, timestamp2, time.Millisecond**deltaMs[0])
	} else {
		assert.WithinDuration(t, timestamp1, timestamp2, time.Millisecond*50)
	}
}

func Assert_DocJson(t *testing.T, actualDoc interface{}, expectedDoc interface{}) {
	actualDocJson, err := json.Marshal(actualDoc)
	assert.Nil(t, err)
	expectedDocJson, err := json.Marshal(expectedDoc)
	assert.Nil(t, err)
	assert.JSONEq(t, string(expectedDocJson), string(actualDocJson))

}
