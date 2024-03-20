package test_service

import (
	"bytes"
	"context"
	"go-app/internals/config"
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"go-app/service"
	"io"
	"net/http"
	"sync"

	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDemoServiceImpl_Account_Create(t *testing.T) {

	tsi := NewTestService(t)
	defer tsi.Clean()

	type fields struct {
		Ctx     context.Context
		Logger  *zerolog.Logger
		Config  *config.DemoServiceConfig
		Service service.Service
	}
	type args struct {
		ctx  context.Context
		opts *schema.Account_CreateOpts
	}

	type TC struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		err       error
		prepare   func(tt *TC)
		validate  func(tt *TC, got interface{})
		Timestamp time.Time
	}

	tests := []TC{
		{
			name: "success",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetTestConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Account_CreateOpts{
					AccountHolderName: gofakeit.Name(),
				},
			},
			wantErr: false,
			prepare: func(tt *TC) {

			},
			validate: func(tt *TC, got interface{}) {
				resp := got.(*schema.Account_CreateResp)
				assert.False(t, resp.ID.IsZero())
				assert.Equal(t, tt.args.opts.AccountHolderName, resp.AccountHolderName)
				assert.NotEmpty(t, resp.UniqueAccountID)
				assert.Len(t, resp.UniqueAccountID, 36)
				assert.Zero(t, resp.Balance)
				// testing for correct timestamp
				allowedDurationDeviation := 3 * time.Second
				Assert_TimestampDuration(t, resp.CreatedAt, tt.Timestamp, &allowedDurationDeviation)

				var doc model.Account
				if err := Get_DocByFilter(
					tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl),
					bson.M{"_id": resp.ID},
					&doc,
				); err != nil {
					assert.Nil(t, err)
				}
				// fmt.Println(doc.CreatedAt, resp.CreatedAt.UTC())
				// testing if document actually exists inside database
				Assert_DocCount(
					t,
					tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl),
					bson.M{"_id": resp.ID},
					1,
				)

				// testing if 2 struct are similar
				assert.Equal(
					t,
					model.Account{
						ID:                resp.ID,
						UniqueAccountID:   resp.UniqueAccountID,
						AccountHolderName: resp.AccountHolderName,
						Balance:           resp.Balance,
						CreatedAt:         resp.CreatedAt,
					},
					doc,
				)

				// testing document value inside database
				Assert_DocJson(
					t,
					doc,
					model.Account{
						ID:                resp.ID,
						UniqueAccountID:   resp.UniqueAccountID,
						AccountHolderName: resp.AccountHolderName,
						Balance:           resp.Balance,
						CreatedAt:         resp.CreatedAt,
					},
				)
			},
		},
		{
			name: "error_account_name_already_exists",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Account_CreateOpts{
					AccountHolderName: gofakeit.Name(),
				},
			},
			wantErr: false,
			prepare: func(tt *TC) {
				m := model.Account{
					ID:                primitive.NewObjectID(),
					UniqueAccountID:   gofakeit.UUID(),
					AccountHolderName: tt.args.opts.AccountHolderName,
					Balance:           0,
					CreatedAt:         tt.Timestamp,
				}
				id, err := tt.fields.Service.MongoDB().
					Cli().
					Database(model.BankDB).
					Collection(model.AccountColl).
					InsertOne(context.TODO(), m)

				fmt.Println(id, err)
			},
			validate: func(tt *TC, got interface{}) {
				resp := got.(*schema.Account_CreateResp)
				fmt.Println(resp.ID)
				Assert_DocCount(
					t,
					tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl),
					bson.M{},
					3, // Note: Here We are specifying 3 DOC as 1 DOC was inserted in previous test.
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Timestamp = fakeUTCNow
			tt.prepare(&tt)
			dsi := &service.DemoServiceImpl{
				Ctx:     tt.fields.Ctx,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
				Service: tt.fields.Service,
			}
			got, err := dsi.Account_Create(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DemoServiceImpl.Account_Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			tt.validate(&tt, got)
		})
	}
}

func TestDemoServiceImpl_Transaction_Create(t *testing.T) {
	tsi := NewTestService(t)
	defer tsi.Clean()

	type fields struct {
		Ctx     context.Context
		Logger  *zerolog.Logger
		Config  *config.DemoServiceConfig
		Service service.Service
	}
	type args struct {
		ctx  context.Context
		opts *schema.Transaction_CreateOpts
	}

	type TC struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		err       error
		prepare   func(tt *TC)
		validate  func(tt *TC)
		Timestamp time.Time
	}

	tests := []TC{
		{
			name: "invalid credit account id",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetTestConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: primitive.NewObjectID(),
					DebitAccountID:  primitive.NewObjectID(),
					Amount:          100,
				},
			},
			wantErr: true,
			err:     errors.New("invalid credit account"),
			prepare: func(tt *TC) {},
			validate: func(tt *TC) {
				tranColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl)
				Assert_DocCount(t, tranColl, bson.M{}, 0)
			},
		},
		{
			name: "0 credit balance insufficient balance",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					DebitAccountID: primitive.NewObjectID(),
					Amount:         100,
				},
			},
			wantErr: true,
			err:     errors.New("insufficient balance"),
			prepare: func(tt *TC) {
				accountColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
				demoAccount := CreateDemoAccountWithZeroBalance(t, accountColl)
				tt.args.opts.CreditAccountID = demoAccount.ID
			},
			validate: func(tt *TC) {
				tranColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl)
				Assert_DocCount(t, tranColl, bson.M{}, 0)
			},
		},
		{
			name: "less credit balance insufficient balance",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					DebitAccountID: primitive.NewObjectID(),
					Amount:         100,
				},
			},
			wantErr: true,
			err:     errors.New("insufficient balance"),
			prepare: func(tt *TC) {
				accountColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
				demoAccount := CreateDemoAccountWithBalance(t, accountColl, 99)
				tt.args.opts.CreditAccountID = demoAccount.ID
			},
			validate: func(tt *TC) {
				tranColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl)
				Assert_DocCount(t, tranColl, bson.M{}, 0)
			},
		},
		{
			name: "no debit account exists",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{

					DebitAccountID: primitive.NewObjectID(),
					Amount:         100,
				},
			},
			wantErr: true,
			err:     errors.New("failed to get debit account: mongo: no documents in result"),
			prepare: func(tt *TC) {
				accountColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
				demoAccount := CreateDemoAccountWithBalance(t, accountColl, 100)
				tt.args.opts.CreditAccountID = demoAccount.ID
			},
			validate: func(tt *TC) {
				// checking if transaction is reported in transaction coll
				tranColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl)
				Assert_DocCount(t, tranColl, bson.M{}, 0)

				// checking if there is an update in credit account
				var account model.Account
				accountColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
				Get_DocByFilter(accountColl, bson.M{"_id": tt.args.opts.CreditAccountID}, &account)
				assert.Equal(t, float32(100), account.Balance)
			},
		},
		{
			name: "success with zero debit account balance",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					Amount: 100,
				},
			},
			wantErr: false,
			prepare: func(tt *TC) {
				accountColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
				demoCreditAccount := CreateDemoAccountWithBalance(t, accountColl, 100)
				tt.args.opts.CreditAccountID = demoCreditAccount.ID

				demoDebitAccount := CreateDemoAccountWithBalance(t, accountColl, 0)
				tt.args.opts.DebitAccountID = demoDebitAccount.ID

			},
			validate: func(tt *TC) {
				// checking if transaction is reported in transaction coll
				tranColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl)
				Assert_DocCount(t, tranColl, bson.M{}, 2)

				// checking if there is an update in credit account
				var creditAccount model.Account
				var debitAccount model.Account
				accountColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
				Get_DocByFilter(accountColl, bson.M{"_id": tt.args.opts.CreditAccountID}, &creditAccount)
				assert.Equal(t, float32(0), creditAccount.Balance)

				Get_DocByFilter(accountColl, bson.M{"_id": tt.args.opts.DebitAccountID}, &debitAccount)
				assert.Equal(t, float32(100), debitAccount.Balance)
				var trans []model.Transaction
				Get_DocsByFilter(tranColl, bson.M{}, &trans)
				assert.Equal(t, 2, len(trans))

				assert.Equal(t, model.CreditTransaction, trans[0].Type)
				assert.Equal(t, trans[0].TransactionID, trans[0].TransactionID)
				assert.Equal(t, tt.args.opts.CreditAccountID, trans[0].CreditAccountID)
				assert.Equal(t, tt.args.opts.DebitAccountID, trans[0].DebitAccountID)
				assert.Equal(t, float32(0), trans[0].ClosingBalance)
				assert.Equal(t, float32(100), trans[0].Amount)

				assert.Equal(t, model.DebitTransaction, trans[1].Type)
				assert.Equal(t, trans[0].TransactionID, trans[1].TransactionID)
				assert.Equal(t, tt.args.opts.CreditAccountID, trans[1].DebitAccountID)
				assert.Equal(t, tt.args.opts.DebitAccountID, trans[1].CreditAccountID)
				assert.Equal(t, float32(100), trans[1].ClosingBalance)
				assert.Equal(t, float32(100), trans[1].Amount)

				Assert_TimestampDuration(t, trans[0].CreatedAt, trans[1].CreatedAt)

			},
		},
		{
			name: "success with already debit account balance",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					Amount: 100,
				},
			},
			wantErr: false,
			prepare: func(tt *TC) {
				accountColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
				demoCreditAccount := CreateDemoAccountWithBalance(t, accountColl, 101)
				tt.args.opts.CreditAccountID = demoCreditAccount.ID

				demoDebitAccount := CreateDemoAccountWithBalance(t, accountColl, 99)
				tt.args.opts.DebitAccountID = demoDebitAccount.ID

			},
			validate: func(tt *TC) {
				// checking if transaction is reported in transaction coll
				tranColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl)
				Assert_DocCount(t, tranColl, bson.M{}, 4)

				// checking if there is an update in credit account
				var creditAccount model.Account
				var debitAccount model.Account
				accountColl := tt.fields.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
				Get_DocByFilter(accountColl, bson.M{"_id": tt.args.opts.CreditAccountID}, &creditAccount)
				assert.Equal(t, float32(1), creditAccount.Balance)

				Get_DocByFilter(accountColl, bson.M{"_id": tt.args.opts.DebitAccountID}, &debitAccount)
				assert.Equal(t, float32(199), debitAccount.Balance)
				var trans []model.Transaction
				Get_DocsByFilter(tranColl, bson.M{}, &trans)
				assert.Equal(t, 4, len(trans))

				// checking credit account
				assert.Equal(t, model.CreditTransaction, trans[2].Type)
				assert.Equal(t, trans[3].TransactionID, trans[2].TransactionID)
				assert.Equal(t, tt.args.opts.CreditAccountID, trans[2].CreditAccountID)
				assert.Equal(t, tt.args.opts.DebitAccountID, trans[2].DebitAccountID)
				assert.Equal(t, float32(1), trans[2].ClosingBalance)
				assert.Equal(t, float32(100), trans[2].Amount)

				// checking debit account
				assert.Equal(t, model.DebitTransaction, trans[3].Type)
				assert.Equal(t, tt.args.opts.CreditAccountID, trans[3].DebitAccountID)
				assert.Equal(t, tt.args.opts.DebitAccountID, trans[3].CreditAccountID)
				assert.Equal(t, float32(199), trans[3].ClosingBalance)
				assert.Equal(t, float32(100), trans[3].Amount)

				Assert_TimestampDuration(t, trans[2].CreatedAt, trans[3].CreatedAt)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Timestamp = fakeUTCNow
			tt.prepare(&tt)
			dsi := &service.DemoServiceImpl{
				Ctx:     tt.fields.Ctx,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
				Service: tt.fields.Service,
			}
			err := dsi.Transaction_Create(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DemoServiceImpl.Transaction_Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			tt.validate(&tt)
		})
	}
}

func TestDemoServiceImpl_Transaction_Create_ACID(t *testing.T) {
	t.Parallel()
	tsi := NewTestService(t)
	defer tsi.Clean()

	type fields struct {
		Ctx     context.Context
		Logger  *zerolog.Logger
		Config  *config.DemoServiceConfig
		Service service.Service
	}
	type args struct {
		ctx  context.Context
		opts *schema.Transaction_CreateOpts
	}

	type TC struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		err       error
		prepare   func(tt *TC)
		validate  func(tt *TC)
		Timestamp time.Time
	}

	accountColl := tsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
	transColl := tsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl)
	demoCreditAccount := CreateDemoAccountWithBalance(t, accountColl, 1000)
	demoDebitAccount := CreateDemoAccountWithBalance(t, accountColl, 0)

	tests := []TC{
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetTestConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
					Amount:          50,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Timestamp = fakeUTCNow
			tt.prepare(&tt)
			dsi := &service.DemoServiceImpl{
				Ctx:     tt.fields.Ctx,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
				Service: tt.fields.Service,
			}
			err := dsi.Transaction_Create(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DemoServiceImpl.Transaction_Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			tt.validate(&tt)
		})
	}

	Assert_DocCount(t, transColl, bson.M{}, 40)

	var trans []model.Transaction
	err := Get_DocsByFilter(transColl, bson.M{}, &trans)
	assert.Nil(t, err)
	assert.Len(t, trans, 40)

	var updatedDemoCreditAccount, updatedDemoDebitAccount model.Account
	err = Get_DocByFilter(accountColl, bson.M{"_id": demoCreditAccount.ID}, &updatedDemoCreditAccount)
	assert.Nil(t, err)

	err = Get_DocByFilter(accountColl, bson.M{"_id": demoDebitAccount.ID}, &updatedDemoDebitAccount)
	assert.Nil(t, err)

	assert.Equal(t, float32(0), updatedDemoCreditAccount.Balance)
	assert.Equal(t, float32(1000), updatedDemoDebitAccount.Balance)
}

func TestDemoServiceImpl_Transaction_Create_ACID_2(t *testing.T) {

	tsi := NewTestService(t)
	defer tsi.Clean()

	type fields struct {
		Ctx     context.Context
		Logger  *zerolog.Logger
		Config  *config.DemoServiceConfig
		Service service.Service
	}
	type args struct {
		ctx  context.Context
		opts *schema.Transaction_CreateOpts
	}

	type TC struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		prepare   func(tt *TC)
		validate  func(tt *TC)
		Timestamp time.Time
	}

	var wg sync.WaitGroup
	max_amount := 1000
	iterations := 100
	transaction_amount := max_amount / iterations
	accountColl := tsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
	transColl := tsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl)
	demoCreditAccount := CreateDemoAccountWithBalance(t, accountColl, float32(max_amount))
	demoDebitAccount := CreateDemoAccountWithBalance(t, accountColl, 0)

	tests := []TC{
		{
			name: "Transaction",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetTestConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.Transaction_CreateOpts{
					CreditAccountID: demoCreditAccount.ID,
					DebitAccountID:  demoDebitAccount.ID,
				},
			},
			wantErr:  false,
			prepare:  func(tt *TC) {},
			validate: func(tt *TC) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Timestamp = fakeUTCNow
			tt.prepare(&tt)
			dsi := &service.DemoServiceImpl{
				Ctx:     tt.fields.Ctx,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
				Service: tt.fields.Service,
			}

			tt.args.opts.Amount = float32(transaction_amount)
			for i := 0; i < iterations; i++ {
				wg.Add(1)
				go func(j int) {
					defer wg.Done()
					err := dsi.Transaction_Create(tt.args.ctx, tt.args.opts)
					assert.Nil(t, err)
				}(i)
			}
		})
	}

	wg.Wait()
	Assert_DocCount(t, transColl, bson.M{}, iterations*2)

	var trans []model.Transaction
	err := Get_DocsByFilter(transColl, bson.M{}, &trans)
	assert.Nil(t, err)
	assert.Len(t, trans, iterations*2)

	var updatedDemoCreditAccount, updatedDemoDebitAccount model.Account
	err = Get_DocByFilter(accountColl, bson.M{"_id": demoCreditAccount.ID}, &updatedDemoCreditAccount)
	assert.Nil(t, err)

	err = Get_DocByFilter(accountColl, bson.M{"_id": demoDebitAccount.ID}, &updatedDemoDebitAccount)
	assert.Nil(t, err)

	assert.Equal(t, float32(0), updatedDemoCreditAccount.Balance)
	assert.Equal(t, float32(max_amount), updatedDemoDebitAccount.Balance)
}

func TestDemoServiceImpl_GetAccountDetailWithTransactions(t *testing.T) {
	tsi := NewTestService(t)
	defer tsi.Clean()

	type fields struct {
		Ctx     context.Context
		Logger  *zerolog.Logger
		Config  *config.DemoServiceConfig
		Service service.Service
	}

	type args struct {
		ctx  context.Context
		opts *schema.AccountTransaction_GetOpts
	}

	type TC struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		err       error
		prepare   func(tt *TC)
		validate  func(tt *TC, got *schema.Account_Get)
		Timestamp time.Time
	}

	tests := []TC{
		{
			name: "no account exists",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetTestConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx: context.TODO(),
				opts: &schema.AccountTransaction_GetOpts{
					ID: primitive.NewObjectID(),
				},
			},
			wantErr: true,
			err:     errors.New("no account found"),
			prepare: func(tt *TC) {},
			validate: func(tt *TC, got *schema.Account_Get) {

			},
		},
		{
			name: "success with 10 transactions",
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			args: args{
				ctx:  context.TODO(),
				opts: &schema.AccountTransaction_GetOpts{},
			},
			wantErr: false,
			prepare: func(tt *TC) {
				accountColl := tsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.AccountColl)
				transColl := tsi.Service.MongoDB().Cli().Database(model.BankDB).Collection(model.TransactionColl)
				creditAccount := CreateDemoAccountWithBalance(t, accountColl, 1000)
				debitAccount := CreateDemoAccountWithZeroBalance(t, accountColl)
				tt.args.opts.ID = creditAccount.ID
				CreateDemoTransactions(t, transColl, *creditAccount, *debitAccount, 10)
			},
			validate: func(tt *TC, got *schema.Account_Get) {
				assert.Equal(t, tt.args.opts.ID, got.ID)
				assert.Len(t, got.Transactions, 20)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Timestamp = fakeUTCNow
			tt.prepare(&tt)
			dsi := &service.DemoServiceImpl{
				Ctx:     tt.fields.Ctx,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
				Service: tt.fields.Service,
			}
			got, err := dsi.GetAccountDetailWithTransactions(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DemoServiceImpl.GetAccountDetailWithTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			tt.validate(&tt, got)
		})
	}
}

func TestDemoServiceImpl_CallAPIForMock(t *testing.T) {

	tsi := NewTestService(t)
	defer tsi.Clean()

	type fields struct {
		Ctx     context.Context
		Logger  *zerolog.Logger
		Config  *config.DemoServiceConfig
		Service service.Service
	}
	type args struct {
		ctx context.Context
		url string
	}

	type TC struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
		prepare func(tt *TC)
	}

	tests := []TC{
		{
			name: "mock success",
			args: args{
				ctx: context.TODO(),
				url: "https://google.com",
			},
			fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetTestConfigFromFile().AppConfig.ServiceConfig.DemoServiceConfig,
				Service: tsi.Service,
			},
			wantErr: false,
			want:    true,
			prepare: func(tt *TC) {
				mockHttpService := tt.fields.Service.GetHTTPService().(*mock.MockHTTP)
				mockHttpService.EXPECT().Get(tt.args.url).Return(
					&http.Response{Status: "200 OK",
						StatusCode:    200,
						Proto:         "HTTP/1.1",
						ProtoMajor:    1,
						ProtoMinor:    1,
						Body:          io.NopCloser(bytes.NewBufferString("Hello world")),
						ContentLength: int64(len("Hello world")),
						Header:        make(http.Header, 0)},
					nil).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsi := &service.DemoServiceImpl{
				Ctx:     tt.fields.Ctx,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
				Service: tt.fields.Service,
			}
			tt.prepare(&tt)
			got, err := dsi.CallAPIForMock(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DemoServiceImpl.CallAPIForMock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DemoServiceImpl.CallAPIForMock() = %v, want %v", got, tt.want)
			}
		})
	}
}
