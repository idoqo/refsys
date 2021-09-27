package db

import "fmt"

type PayoutType string
type PayoutStatus string

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

var 	ErrNoRow = fmt.Errorf("no matching row found")

type User struct {
	ID int64 `json:"id,omitempty"`
	Username string `json:"username"`
	ReferralCode string `json:"referral_code,omitempty"`
	Referrer string `json:"referrer,omitempty"`
}

type Payout struct {
	ID int64 `json:"id,omitempty"`
	UserID int64 `json:"user_id"`
	CheckpointID int64 `json:"checkpoint_id"`
	Amount int64 `json:"amount"`
	Status PayoutStatus `json:"status"`
	Type PayoutType `json:"type"`
}

type Database interface
{
	SaveUser(username, hashedPassword, referrer string) (id int64, referralCode string, err error)
	UserByRefCode(refCode string) (User, error)
	CheckAndTriggerPayout(referrer User, payoutType PayoutType) (payout Payout, shouldTrigger bool, err error)
}
