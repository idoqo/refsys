package main

import (
	"github.com/dimfeld/httptreemux/v5"
	"github.com/rs/zerolog"
	"gitlab.com/idoko/refsys/db"
	api "gitlab.com/idoko/refsys/pkg/api"
	"net/http"
	"os"
	"strconv"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	dbCfg, err := getDbConfig()
	if err != nil {
		logger.Err(err).Msg("configuring database")
		os.Exit(1)
	}

	pg, err := db.InitPostgres(dbCfg, logger)
	if err != nil {
		logger.Err(err).Msg("database connection")
		os.Exit(1)
	}

	router := httptreemux.NewContextMux()
	rg := router.NewGroup("/api")

	server := api.NewApi(pg, logger)
	server.RegisterRouter(rg)

	logger.Err(http.ListenAndServe(":8000", router)).Msg("starting server")
}

func getDbConfig() (db.Config, error){
	var cfg db.Config
	var err error
	var port int
	if port, err = strconv.Atoi(os.Getenv("POSTGRES_DB_PORT")); err != nil {
		return cfg, err
	}

	cfg = db.Config{
		Host: os.Getenv("POSTGRES_DB_HOST"),
		Port: port,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Schema: os.Getenv("POSTGRES_DB"),
	}
	return cfg, nil
}