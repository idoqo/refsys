package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"math/rand"
	"time"
)

type Config struct {
	Host string
	Username string
	Password string
	Port int
	Schema string

}
type PgDb struct {
	Conn *sql.DB
	Logger zerolog.Logger
}

func InitPostgres(cfg Config, logger zerolog.Logger) (PgDb, error){
	db := PgDb{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Schema)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	db.Logger = logger
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

func (p PgDb) SaveUser(username, hashedPassword, referrer string) (id int64, referralCode string, err error) {
	refCode := generateRefCode(6)

	if referrer == "" {
		insertQuery := "INSERT INTO users(username, password, referral_code) VALUES ($1, $2, $3) RETURNING id"
		err = p.Conn.QueryRow(insertQuery, username, hashedPassword, refCode).Scan(&id)
	} else {
		insertQuery := "INSERT INTO users(username, password, referrer, referral_code) VALUES ($1, $2, $3, $4) RETURNING id"
		err = p.Conn.QueryRow(insertQuery, username, hashedPassword, referrer, refCode).Scan(&id)
	}

	return id, referralCode, err
}

func (p PgDb) UserByRefCode(refCode string) (User, error) {
	var u User
	query := "SELECT id, username, referral_code FROM users WHERE referral_code=$1 LIMIT 1"
	err := p.Conn.QueryRow(query, refCode).Scan(&u.ID, &u.Username, &u.ReferralCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNoRow
		}
		return u, err
	}
	return u, nil
}

// CheckAndTriggerPayout  it trigger payouts by doing the following:
// - find the last time we paid out this user (if it exists), and the id of the user(or friend) that triggered this payout - this friend is the last checkpoint
// - select all friends that come after the checkpoint (e.g for signups, it will check all users with the same referrer)
// - if the matching friends are up to three, it records the 3rd person as a checkpoint
func (p PgDb) CheckAndTriggerPayout(referrer User, payoutType PayoutType) (Payout, bool, error) {
	var lastCheckpoint int
	var payout Payout

	checkpointQ := "SELECT checkpoint_id FROM payouts WHERE user_id=$1 AND activity_type=$2 ORDER BY id DESC LIMIT 1"
	err := p.Conn.QueryRow(checkpointQ, referrer.ID, payoutType).Scan(&lastCheckpoint)
	if err != nil && err != sql.ErrNoRows {
		return payout, false, err
	}

	friends := make([]User, MinReferrals)
	if payoutType == Signups {
		err = p.copyFriendsForSignups(friends, referrer, lastCheckpoint)
	} else if payoutType == Transactions {
		err = p.copyFriendsForTransactions(friends, referrer, lastCheckpoint)
	}
	if err != nil {
		return payout, false, err
	}

	p.Logger.Info().Msgf("processing payouts because %+v", friends)

	// pick the last user in the list to use as checkpoint
	newCheckpoint := friends[MinReferrals-1]
	if newCheckpoint.ID == 0 {
		return payout, false, nil
	}

	var amount string
	newPayoutQ := "INSERT INTO payouts(user_id, activity_type, checkpoint_id, amount, status) VALUES($1, $2, $3, $4, $5) RETURNING id, user_id, activity_type, checkpoint_id, amount, status"
	err = p.Conn.QueryRow(newPayoutQ, referrer.ID, Signups, newCheckpoint.ID, ReferralBonus, Pending).Scan(
		&payout.ID, &payout.UserID, &payout.Type, &payout.CheckpointID, &amount, &payout.Status,
	)
	if err != nil {
		return payout, false, err
	} else {
		payout.Amount, _ = decimal.NewFromString(amount)
		payout.Username = referrer.Username
		return payout, true, nil
	}
}

func (p PgDb) copyFriendsForSignups(friends []User, referrer User, lastCheckpoint int) error {
	// checkpoint remains 0 if no payout has happened for the user, so it's safe to use in querying friends.
	friendsQ := "SELECT id, username, referrer FROM users WHERE id > $1 AND referrer = $2 ORDER BY id LIMIT $3"
	rows, err := p.Conn.Query(friendsQ, lastCheckpoint, referrer.ReferralCode, MinReferrals)
	if err != nil {
		return err
	}

	i := 0
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Referrer)
		if err != nil {
			return err
		}
		friends[i] = u
		i++
	}
	return nil
}

func (p PgDb) copyFriendsForTransactions(friends []User, referrer User, lastCheckpoint int) error {
	// select all friends (user_id where referrer = referrer.ReferralCode)
	// select from transactions where sender_id in [friends_id] and amount > 200 and id > lastcheckpoint

	friendsQ := "SELECT id, username, referrer FROM users WHERE id > $1 AND referrer = $2 ORDER BY id LIMIT $3"
	rows, err := p.Conn.Query(friendsQ, lastCheckpoint, referrer.ReferralCode, MinReferrals)
	if err != nil {
		return err
	}

	i := 0
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Referrer)
		if err != nil {
			return err
		}
		friends[i] = u
		i++
	}
	return nil
}

func generateRefCode(length int) string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
