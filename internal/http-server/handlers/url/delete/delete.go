package delete

import (
	"errors"
	resp "mod_shortener/internal/lib/api/responce"
	"mod_shortener/internal/lib/logger/sl"
	"mod_shortener/internal/storage"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(id int64) (int64, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.url.delete"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Info(op, slog.String("id", "id is empty"))

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		idUrl, err := strconv.Atoi(id)
		if err != nil {
			log.Error(op, sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		res, err := urlDeleter.DeleteURL(int64(idUrl))
		if errors.Is(err, storage.ErrElemNotFount) {
			log.Info("id not found", "id", id)

			render.JSON(w, r, "id not found")

			return
		}

		if err != nil {
			log.Error(op, sl.Err(err))

			render.JSON(w, r, "internal error")

			return
		}

		if res < 1 {
			log.Info(op, "res", "not deleted")

			render.JSON(w, r, resp.Error("not delete element"))
		}

		log.Info("url delete", "id", id)

		render.JSON(w, r, resp.Ok())

		return
	}
}
