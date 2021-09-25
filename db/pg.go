package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
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