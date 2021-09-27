package db

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type PayoutType string
type PayoutStatus string

type TransactionType string

const (
	MinReferrals = 3
	ReferralBonus = 50
)

const (
	Signups PayoutType = "signups"
	Transactions = "transactions"
)

const (
	Pending PayoutStatus = "pending"
	Paid = "paid"
)

const (
	Debit TransactionType = "debit"
	Credit = "credit"
)

var (
	ErrNoRow = fmt.Errorf("no matching row found")
	ErrInsufficientFunds = fmt.Errorf("insufficient balance")
)

type User struct {
	ID int64 `json:"id,omitempty"`
	Username string `json:"username"`
	ReferralCode string `json:"referral_code,omitempty"`
	Referrer string `json:"referrer,omitempty"`
}

type Payout struct {
	ID int64 `json:"id,omitempty"`
	UserID int64 `json:"user_id"`
	Username string `json:"username,omitempty"`
	CheckpointID int64 `json:"checkpoint_id"`
	Amount decimal.Decimal `json:"amount"`
	Status PayoutStatus `json:"status"`
	Type PayoutType `json:"type"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Transaction struct {
	ID int64 `json:"id"`
	Reference string `json:"reference"`
	SenderID int64 `json:"sender_id"`
	RecipientID int64 `json:"recipient_id"`
	Amount decimal.Decimal `json:"amount"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Database interface
{
	SaveUser(username, hashedPassword, referrer string) (id int64, referralCode string, err error)
	UserByRefCode(refCode string) (User, error)
	UsersReferrer(userId int64) (User, error)

	SaveTransaction(transaction *Transaction) error

	CheckAndTriggerPayout(referrer User, payoutType PayoutType) (payout Payout, shouldTrigger bool, err error)
}
