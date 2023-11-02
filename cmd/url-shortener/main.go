package main

import (
	"log/slog"
	"mod_shortener/internal/config"
	"mod_shortener/internal/http-server/handlers/url/delete"
	"mod_shortener/internal/http-server/handlers/url/redirect"
	"mod_shortener/internal/http-server/handlers/url/save"
	"mod_shortener/internal/http-server/handlers/users"
	"mod_shortener/internal/lib/logger/sl"
	"mod_shortener/internal/lib/my_middleware"
	"mod_shortener/internal/storage/sqlite"
	"net/http"
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
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	storage, err := sqlite.New(cfg.Storage)
	
	if err != nil {
		log.Error("failed init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
			"aaa":               "aaa",
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{id}", delete.New(log, storage))
	})

	router.Post("/users/reg", users.New(log, storage))
	router.Post("/users/auth", users.Auth(log, storage, cfg))
	router.Post("/users/refresh", users.Refresh(log, storage, cfg))

	router.Group(func(r chi.Router) {
		r.Use(my_middleware.JWT)
		r.Get("/{alias}", redirect.New(log, storage))
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
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
