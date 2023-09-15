package main

import (
	"log/slog"
	"mod_shortener/internal/config"
	"mod_shortener/internal/lib/logger/sl"
	"mod_shortener/internal/storage/sqlite"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// todo: init config: cleanenv
	cfg := config.MustLoad()

	// todo: init logger: slog
	log := setupLogger(cfg.Env)

	// log = log.With(slog.String("env006", cfg.Env))
	// log.Info("starting url-shortener", slog.String("env007", cfg.Env))
	// log.Debug("this is Debug")

	// todo: init storage: sqlite
	storage, err := sqlite.New(cfg.Storage)
	if err != nil {
		log.Error("failed init storage", sl.Err(err))
		os.Exit(1)
	}

	res, err := storage.GetURL("ya")
	if err != nil {
		log.Error("res", sl.Err(err))
	}

	log.Info(res)

	/* id, err := storage.SaveURL("www.google.ru", "gl")
	if err != nil {
		log.Error("err", sl.Err(err))
		os.Exit(1)
	}

	log.Info("save url", slog.Int64("id", id))

	resDel, err := storage.DeleteURL(id)

	log.Info("delete by id", slog.Int64("Удалено строк:", resDel))

	if err != nil {
		log.Error("delElem", sl.Err(err))
	} */

	// todo: init router: chi, "chi render"
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// todo: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
