package users

import (
	"log/slog"
	resp "mod_shortener/internal/lib/api/responce"
	"mod_shortener/internal/lib/api/user"
	"mod_shortener/internal/lib/logger/sl"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Responce struct {
	resp.Responce
	Id int64 `json:"id,omitempty"`
}

type UserStorage interface {
	AddUser(user *user.User, log *slog.Logger) (int64, error)
	GetPass()
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

func Auth(log *slog.Logger, storage UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "users.Auth"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		user := &user.User{}
		render.DecodeJSON(r.Body, user)
		// получаем логин

		storage.GetPass()
		// сравниваем
		// вызываем или возвращаем управление jwt

	}
}
