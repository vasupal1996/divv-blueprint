package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertOneOpts struct {
	Name string `json:"name" validate:"required,ne=LoremIpsumLoremIpsum,max=12"`
}

type Account_CreateOpts struct {
	AccountHolderName string `json:"account_holder_name"`
}

type Account_CreateResp struct {
	ID                primitive.ObjectID `json:"id"`
	UniqueAccountID   string             `json:"account_id,omitempty"`
	AccountHolderName string             `json:"account_holder_name"`
	Balance           float32            `json:"balance"`
	CreatedAt         time.Time          `json:"created_at"`
}

type Transaction_CreateOpts struct {
	TransactionID   string
	CreditAccountID primitive.ObjectID `json:"credit_account_id"`
	DebitAccountID  primitive.ObjectID `json:"debit_account_id"`
	Amount          float32            `json:"amount,omitempty"`
}

type AccountTransaction_GetOpts struct {
	ID primitive.ObjectID `json:"id"`
}

type Account_Get struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	UniqueAccountID   string             `json:"account_id,omitempty" bson:"account_id,omitempty"`
	AccountHolderName string             `json:"account_holder_name" bson:"account_holder_name"`
	Balance           float32            `json:"balance" bson:"balance"`
	Transactions      []Transaction_Get  `json:"transactions" bson:"transactions"`
}

type Transaction_Get struct {
	TransactionID   string             `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	CreditAccountID primitive.ObjectID `json:"credit_account_id,omitempty" bson:"credit_account_id,omitempty"`
	DebitAccountID  primitive.ObjectID `json:"debit_account_id,omitempty" bson:"debit_account_id,omitempty"`
	Type            string             `json:"type,omitempty" bson:"type,omitempty"`
	Amount          float32            `json:"amount,omitempty" bson:"amount,omitempty"`
	ClosingBalance  float32            `json:"closing_balance,omitempty" bson:"closing_balance,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
