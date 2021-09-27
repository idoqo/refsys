package api

import (
	"encoding/json"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/rs/zerolog"
	"gitlab.com/idoko/refsys/db"
	"gitlab.com/idoko/refsys/pkg/service"
	"net/http"
)

type Api struct {
	Logger zerolog.Logger
	Service service.Service
}

func NewApi(database db.PgDb, logger zerolog.Logger) Api {
	svc := service.Service {
		Db: database,
		Logger: logger,
	}
	return Api{
		Logger: logger,
		Service: svc,
	}
}

// returns an HTTP response in JSON format
func (api Api) res(statusCode int, w http.ResponseWriter, body interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(body)
}

func (api Api) RegisterRouter(rg *httptreemux.ContextGroup) {
	rg.GET("/healthz", api.HealthCheck)
	rg.POST("/register", api.Signup)
	rg.POST("/transaction", api.Transaction)
}

func (api Api) HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := api.res(http.StatusOK, w, map[string]string{
		"status": "alive",
	})
	if err != nil {
		api.Logger.Err(err).Msg("failed sending response")
	}
	return
}

type userReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Referrer string `json:"referrer"`
}

func (api Api) Signup(w http.ResponseWriter, r *http.Request) {
	var ur userReq
	jd := json.NewDecoder(r.Body)
	if err := jd.Decode(&ur); err != nil {
		api.Logger.Err(err).Msg("parsing user request")
		_ = api.res(http.StatusBadRequest, w, map[string]string{
			"error": "failed to parse request",
		})
		return
	}

	if ur.Username == "" || ur.Password == "" {
		_ = api.res(http.StatusBadRequest, w, map[string]string{
			"error": "username and/or password cannot be blank",
		})
		return
	}

	u, err := api.Service.CreateUser(ur.Username, ur.Password, ur.Referrer)
	if err != nil {
		api.Logger.Err(err).Msg("sign up")
		_ = api.res(http.StatusInternalServerError, w, map[string]string{
			"error": "could not complete signup",
		})
		return
	}

	_ = api.res(http.StatusCreated, w, u)
}

func (api Api) Transaction(w http.ResponseWriter, r *http.Request) {

}
