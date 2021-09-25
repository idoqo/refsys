package api

import (
	"github.com/dimfeld/httptreemux/v5"
	"github.com/rs/zerolog"
	"gitlab.com/idoko/refsys/db"
)

type Api struct {
	Db db.PgDb
	Logger zerolog.Logger
}

func NewApi(database db.PgDb, logger zerolog.Logger) Api {
	return Api{
		Db: database,
		Logger: logger,
	}
}

func (api Api) RegisterRouter(rg *httptreemux.ContextGroup) {

}
