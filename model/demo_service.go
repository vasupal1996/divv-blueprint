package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DemoServiceDB   = "demo_service"
	DemoServiceColl = "demo_service"

	BankDB          = "demo_bank"
	AccountColl     = "account"
	TransactionColl = "transaction"
)

const (
	CreditTransaction = "credit"
	DebitTransaction  = "debit"
)

type DemoService struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Account struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UniqueAccountID   string             `json:"account_id,omitempty" bson:"account_id,omitempty"`
	AccountHolderName string             `json:"account_holder_name,omitempty" bson:"account_holder_name,omitempty"`
	Balance           float32            `json:"balance,omitempty" bson:"balance,omitempty"`
	CreatedAt         time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Transaction struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TransactionID   string             `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	CreditAccountID primitive.ObjectID `json:"credit_account_id,omitempty" bson:"credit_account_id,omitempty"`
	DebitAccountID  primitive.ObjectID `json:"debit_account_id,omitempty" bson:"debit_account_id,omitempty"`
	Type            string             `json:"type,omitempty" bson:"type,omitempty"`
	Amount          float32            `json:"amount,omitempty" bson:"amount,omitempty"`
	ClosingBalance  float32            `json:"closing_balance,omitempty" bson:"closing_balance,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
