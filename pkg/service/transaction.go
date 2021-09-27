package service

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.com/idoko/refsys/db"
)

func (s Service) CreateTransaction(senderId, recipientId int, amount decimal.Decimal, description string) (db.Transaction, error) {
	t := db.Transaction {
		SenderID: int64(senderId),
		RecipientID: int64(recipientId),
		Amount: amount,
		Description: description,
	}

	err := s.Db.SaveTransaction(&t)
	if err != nil {
		return t, err
	}

	minAmountForReferral := decimal.NewFromInt(db.MinAmountForReferral)
	if t.Amount.Cmp(minAmountForReferral) >= 0 {
		go s.TransactionCreatedEvent(t)
	}
	return t, nil
}

func (s Service) TransactionCreatedEvent(t db.Transaction) {
	referrer, err := s.Db.UsersReferrer(t.SenderID)
	if err != nil && err != db.ErrNoRow{
		s.Logger.Err(err).Msg("getting user's referrer")
		return
	}

	payout, shouldTrigger, err := s.Db.CheckAndTriggerPayout(referrer, db.Transactions)
	if err != nil {
		s.Logger.Err(err).Msg("checking payout for transaction")
		return
	}

	if shouldTrigger {
		s.Logger.Info().Msg("payout due")
		s.Logger.Info().Msg(fmt.Sprintf("%+v", payout))
	}
}