package save

import (
	resp "mod_shortener/internal/lib/api/responce"
	"mod_shortener/internal/lib/logger/sl"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias, omitempty"`
}

type Responce struct {
	resp.Responce
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handler.url.save.new"

		log = log.With(
			slog.String("op", op),
			slog.With("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
	}
}
