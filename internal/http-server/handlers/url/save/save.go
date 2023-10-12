package save

import (
	"errors"
	"log/slog"
	resp "mod_shortener/internal/lib/api/responce"
	"mod_shortener/internal/lib/logger/sl"
	"mod_shortener/internal/lib/random"
	"mod_shortener/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Responce struct {
	resp.Responce
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handler.url.save.new"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// валидируем структуру req
		if err := validator.New().Struct(req); err != nil {
			validatError := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validatError))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.GetRandomByLength(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("URL already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("URL already exists"))
			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("URL added", slog.Int64("id", id))

		responceOk(w, r, alias)
	}
}

func responceOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Responce{
		Responce: resp.Ok(),
		Alias:    alias,
	})
}
