//go:generate $GOPATH/bin/mockgen -destination=../mock/mock_demo_service.go -package=mock go-app/service DemoService
package service

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"io"

	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DemoService interface {
	DemoFunc(ctx context.Context) string
	SentryDemoFunc(ctx context.Context) string
	InsertOne(ctx context.Context, opts *schema.InsertOneOpts) (primitive.ObjectID, error)
	CallAPIForMock(ctx context.Context, url string) (bool, error)

	// MongoDB Related Operations
	Account_Create(ctx context.Context, opts *schema.Account_CreateOpts) (*schema.Account_CreateResp, error)
	Transaction_Create(ctx context.Context, opts *schema.Transaction_CreateOpts) error

	GetAccountDetailWithTransactions(ctx context.Context, opts *schema.AccountTransaction_GetOpts) (*schema.Account_Get, error)
}

func (dsi *DemoServiceImpl) DemoFunc(ctx context.Context) string {
	dsi.Logger.Info().Ctx(ctx).Msg(dsi.Config.SomeAdditionalData)
	return dsi.Config.SomeAdditionalData
}

func (dsi *DemoServiceImpl) SentryDemoFunc(ctx context.Context) string {
	dsi.Logger.Info().Ctx(ctx).Msg("test log")
	dsi.Logger.Warn().Ctx(context.WithValue(ctx, schema.SentryExtraCtx, map[string]string{"some": "thing"})).Msg(dsi.Config.SomeAdditionalData)
	return dsi.Config.SomeAdditionalData
}

func (dsi *DemoServiceImpl) InsertOne(ctx context.Context, opts *schema.InsertOneOpts) (primitive.ObjectID, error) {
	m := model.DemoService{
		ID:        primitive.NewObjectID(),
		Name:      opts.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	resp, err := dsi.Service.MongoDB().Cli().Database(model.DemoServiceDB).Collection(model.DemoServiceColl).InsertOne(ctx, m)

	if err != nil {
		return primitive.NilObjectID, errors.Wrap(err, "failed to insert document")
	}
	return resp.(primitive.ObjectID), nil
}

func (dsi *DemoServiceImpl) Account_Create(ctx context.Context, opts *schema.Account_CreateOpts) (*schema.Account_CreateResp, error) {
	m := model.Account{
		ID:                primitive.NewObjectID(),
		UniqueAccountID:   uuid.New().String(),
		AccountHolderName: opts.AccountHolderName,
		Balance:           0,
		CreatedAt:         UTCNow(),
	}

	res, err := dsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl).InsertOne(ctx, m)
	if err != nil {
		dsi.Logger.Err(err).Ctx(ctx).Interface("m", m).Msg(err.Error())
		return nil, errors.Wrap(err, "failed to create account")
	}

	resp := schema.Account_CreateResp{
		ID:                res.(primitive.ObjectID),
		UniqueAccountID:   m.UniqueAccountID,
		AccountHolderName: m.AccountHolderName,
		Balance:           m.Balance,
		CreatedAt:         m.CreatedAt,
	}

	return &resp, nil
}

func (dsi *DemoServiceImpl) Transaction_Create(ctx context.Context, opts *schema.Transaction_CreateOpts) error {
	opts.TransactionID = uuid.New().String()
	err := dsi.registerTransaction(ctx, opts)
	return err
}

func (dsi *DemoServiceImpl) registerTransaction(ctx context.Context, opts *schema.Transaction_CreateOpts) error {

	session, err := dsi.Service.MongoDB().Cli().StartSession()
	if err != nil {
		dsi.Logger.Err(err).Msg("failed to start session for transaction")
		return errors.Wrap(err, "failed to start session for transaction")
	}

	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {

		var creditAccount model.Account
		res := dsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl).FindOne(
			sessionContext,
			bson.M{"_id": opts.CreditAccountID},
		)

		if err := res.Decode(&creditAccount); err != nil {
			dsi.Logger.Err(err).Ctx(ctx).Interface("opts", opts).Msg("failed to credit get account")
			return nil, errors.New("invalid credit account")
		}

		if creditAccount.Balance-opts.Amount < 0 {
			return nil, errors.New("insufficient balance")
		}

		// creating transaction
		tc := model.Transaction{
			ID:              primitive.NewObjectID(),
			TransactionID:   opts.TransactionID,
			CreditAccountID: opts.CreditAccountID,
			DebitAccountID:  opts.DebitAccountID,
			Type:            model.CreditTransaction,
			Amount:          opts.Amount,
			ClosingBalance:  creditAccount.Balance - opts.Amount,
			CreatedAt:       UTCNow(),
		}
		_, err := dsi.Service.MongoDB().
			Cli().
			Database(model.BankDB).
			Collection(model.TransactionColl).
			InsertOne(sessionContext, tc)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create transaction")
		}

		// updating balance
		tcUpdate := bson.M{
			"$set": bson.M{
				"balance":    tc.ClosingBalance,
				"updated_at": UTCNow(),
			},
		}

		if _, err := dsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl).UpdateOne(sessionContext, bson.M{"_id": creditAccount.ID}, tcUpdate); err != nil {
			return nil, errors.Wrap(err, "failed to update account balance")
		}

		var debitAccount model.Account
		res = dsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl).FindOne(
			sessionContext,
			bson.M{"_id": opts.DebitAccountID},
		)
		if err := res.Decode(&debitAccount); err != nil {
			dsi.Logger.Err(err).Ctx(ctx).Interface("opts", opts).Msg("failed to get account")
			return nil, errors.Wrap(err, "failed to get debit account")
		}

		// creating transaction
		td := model.Transaction{
			ID:              primitive.NewObjectID(),
			TransactionID:   opts.TransactionID,
			CreditAccountID: opts.DebitAccountID,
			DebitAccountID:  opts.CreditAccountID,
			Type:            model.DebitTransaction,
			Amount:          opts.Amount,
			ClosingBalance:  debitAccount.Balance + opts.Amount,
			CreatedAt:       UTCNow(),
		}
		_, err = dsi.Service.MongoDB().
			Cli().
			Database(model.BankDB).
			Collection(model.TransactionColl).
			InsertOne(sessionContext, td)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create transaction")
		}

		// updating balance
		tdUpdate := bson.M{
			"$set": bson.M{
				"balance":    td.ClosingBalance,
				"updated_at": UTCNow(),
			},
		}

		if _, err := dsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl).UpdateOne(sessionContext, bson.M{"_id": debitAccount.ID}, tdUpdate); err != nil {
			return nil, errors.Wrap(err, "failed to update account balance")
		}

		return nil, nil
	})

	return err
}

func (dsi *DemoServiceImpl) GetAccountDetailWithTransactions(ctx context.Context, opts *schema.AccountTransaction_GetOpts) (*schema.Account_Get, error) {

	var accountResp schema.Account_Get
	res := dsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl).FindOne(ctx, bson.M{"_id": opts.ID})

	if err := res.Decode(&accountResp); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no account found")
		}
		return nil, errors.New("failed to get account")
	}

	var transactions []schema.Transaction_Get

	cur, err := dsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl).Find(ctx, bson.M{
		"$or": bson.A{
			bson.M{
				"credit_account_id": opts.ID,
			},
			bson.M{
				"debit_account_id": opts.ID,
			},
		},
	})

	if err != nil {
		return nil, errors.New("failed to prepare query")
	}

	if err := cur.All(ctx, &transactions); err != nil {
		return nil, errors.New("failed to get transactions")
	}

	accountResp.Transactions = transactions

	return &accountResp, nil
}

func (dsi *DemoServiceImpl) CallAPIForMock(ctx context.Context, url string) (bool, error) {
	print(dsi.DemoFunc(ctx))

	resp, err := dsi.Service.GetHTTPService().Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	return true, nil
}
