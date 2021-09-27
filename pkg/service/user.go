package service

import (
	"fmt"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/idoko/refsys/db"
)

type Service struct {
	Db db.Database
	Logger zerolog.Logger
}

func (s Service) CreateUser(username, plainPassword, referrer string) (db.User, error) {
	u := db.User{}
	password, err := hashPassword(plainPassword)
	if err != nil {
		return u, err
	}

	id, refCode, err := s.Db.SaveUser(username, password, referrer)
	if err != nil {
		return u, err
	}
	u = db.User{
		ID: id,
		Username: username,
		Referrer: referrer,
		ReferralCode: refCode,
	}

	go s.UserCreatedEvent(u)
	return u, nil
}

func hashPassword(plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 14)
	return string(bytes), err
}

func (s Service) UserCreatedEvent(user db.User) {
	if user.Referrer == "" {
		return
	}

	referrer, err := s.Db.UserByRefCode(user.Referrer)
	if err != nil {
		s.Logger.Err(err).Msg("getting referrer details")
		return
	}
	payout, shouldTrigger, err := s.Db.CheckAndTriggerPayout(referrer, db.Signups)
	if err != nil {
		s.Logger.Err(err).Msg("checking payouts")
		return
	}

	if shouldTrigger {
		// maybe we could send a POST request to the notifier endpoint using payout as body
		// for now we just log
		s.Logger.Info().Msg("payout due")
		s.Logger.Info().Msg(fmt.Sprintf("%+v", payout))
	}
}
