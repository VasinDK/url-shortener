package users

import (
	"encoding/json"
	"log/slog"
	"mod_shortener/internal/config"
	resp "mod_shortener/internal/lib/api/responce"
	"mod_shortener/internal/lib/api/user"
	"mod_shortener/internal/lib/crypto"
	"mod_shortener/internal/lib/logger/sl"
	"mod_shortener/internal/lib/random"
	"mod_shortener/internal/lib/token/jwt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Responce struct {
	resp.Responce
	Id int64 `json:"id,omitempty"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserStorage interface {
	AddUser(*user.User, *slog.Logger) (int64, error)
	GetPass(string, *slog.Logger) (string, string, string, error)
	GetUser(string, *slog.Logger) (*user.User, error)
	UpdateRefreshToken(string, time.Time, *slog.Logger, string) (string, error)
}

func New(log *slog.Logger, storage UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "user.add.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var user = &user.User{}

		err := render.DecodeJSON(r.Body, user)

		if err != nil {
			log.Error("failed to decode r.Body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode r.Body"))
			return
		}

		id, err := storage.AddUser(user, log)

		if err != nil {
			log.Error("failed to add user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add user"))
			return
		}

		log.Info("user added", slog.Int64("id", id))

		responseOk(w, r, id)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, id int64) {
	render.JSON(w, r, Responce{
		Responce: resp.Ok(),
		Id:       id,
	})
}

func Auth(log *slog.Logger, storage UserStorage, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "users.Auth"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		user := &user.User{}
		render.DecodeJSON(r.Body, user)

		hash, _, id, err := storage.GetPass(user.Login, log)

		if err != nil {
			log.Error(op+"GetPass", sl.Err(err))
		}

		user.Id = id

		auth := crypto.CheckPassHash(user.Pass, hash)

		if !auth {
			log.Info(http.StatusText(http.StatusUnauthorized))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		setAccessRefreshToken(cfg, user, w, r, log, storage)

		return
	}
}

func Refresh(log *slog.Logger, storage UserStorage, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "user.Refresh"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var request user.User
		var token string

		err := json.NewDecoder(r.Body).Decode(&request)

		_, token, _, err = storage.GetPass(request.Login, log)
		userCurrent, err := storage.GetUser(request.Login, log)

		if err != nil {
			log.Error(op, sl.Err(err))
		}

		if request.Refresh_token == token {
			setAccessRefreshToken(cfg, userCurrent, w, r, log, storage)
		}

		return
	}
}

func setAccessRefreshToken(
	cfg *config.Config,
	user *user.User,
	w http.ResponseWriter,
	r *http.Request,
	log *slog.Logger,
	storage UserStorage,
) {
	const op = "user.setAccessRefreshToken"
	accessToken, err := jwt.CreateJWT(cfg.AccessTokenTTL, cfg.KeyToken, user)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	refreshToken := random.GetRandomByLength(10)

	if err != nil {
		log.Error(op, sl.Err(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	expiresAT := time.Now().Add(time.Hour * time.Duration(cfg.RefreshTokenTTL))

	_, err = storage.UpdateRefreshToken(refreshToken, expiresAT, log, user.Id)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
