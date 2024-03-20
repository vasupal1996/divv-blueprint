package service_test

import (
	"context"
	"go-app/internals/mongodb"
	"go-app/model"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var fakeUTCNow = time.Now().UTC().Truncate(time.Millisecond)

// var fakeNow = time.Now().UTC().Truncate(time.Millisecond)

func CreateDemoAccountWithZeroBalance(t *testing.T, coll mongodb.Collection) *model.Account {
	m := model.Account{
		UniqueAccountID:   gofakeit.UUID(),
		AccountHolderName: gofakeit.Name(),
		Balance:           0,
		CreatedAt:         fakeUTCNow,
	}
	resp, err := coll.InsertOne(context.TODO(), m)
	assert.Nil(t, err)

	var demoAccount model.Account
	if res := coll.FindOne(context.TODO(), bson.M{"_id": resp}); res != nil {
		assert.NotNil(t, res)
		if err := res.Decode(&demoAccount); err != nil {
			assert.Nil(t, err)
		}
	}

	return &demoAccount
}

func CreateDemoAccountWithBalance(t *testing.T, coll mongodb.Collection, bal float32) *model.Account {
	m := model.Account{
		UniqueAccountID:   gofakeit.UUID(),
		AccountHolderName: gofakeit.Name(),
		Balance:           bal,
		CreatedAt:         fakeUTCNow,
	}
	resp, err := coll.InsertOne(context.TODO(), m)
	assert.Nil(t, err)

	var demoAccount model.Account
	if res := coll.FindOne(context.TODO(), bson.M{"_id": resp}); res != nil {
		assert.NotNil(t, res)
		if err := res.Decode(&demoAccount); err != nil {
			assert.Nil(t, err)
		}
	}

	return &demoAccount
}

func CreateDemoTransactions(t *testing.T, coll mongodb.Collection, creditAccount, debitAccount model.Account, transactionCount int) []model.Transaction {
	var transactions []interface{}

	transactionAmount := creditAccount.Balance / float32(transactionCount)

	for i := 0; i < transactionCount; i++ {
		transactionID := gofakeit.UUID()
		ct := model.Transaction{
			TransactionID:   transactionID,
			CreditAccountID: creditAccount.ID,
			DebitAccountID:  debitAccount.ID,
			Type:            model.CreditTransaction,
			Amount:          transactionAmount,
			ClosingBalance:  creditAccount.Balance - transactionAmount,
			CreatedAt:       fakeUTCNow,
		}

		dt := model.Transaction{
			TransactionID:   transactionID,
			CreditAccountID: creditAccount.ID,
			DebitAccountID:  debitAccount.ID,
			Type:            model.DebitTransaction,
			Amount:          transactionAmount,
			ClosingBalance:  debitAccount.Balance + transactionAmount,
			CreatedAt:       fakeUTCNow,
		}

		transactions = append(transactions, ct, dt)
	}

	ids, err := coll.InsertMany(context.TODO(), transactions)
	assert.Len(t, ids, transactionCount*2)
	assert.Nil(t, err)

	var resp []model.Transaction

	for _, j := range transactions {
		resp = append(resp, j.(model.Transaction))
	}

	return resp
}
