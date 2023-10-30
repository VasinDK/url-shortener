package my_middleware

import (
	"context"

	"log/slog"
	"mod_shortener/internal/config"
	"mod_shortener/internal/lib/logger/sl"
	"mod_shortener/internal/lib/token/jwt"
	"net/http"
	"strings"
)

func JWT(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		const op = "my_middleware.JWT"

		cfg := config.MustLoad()

		reqToken := r.Header.Get("Authorization")
		
		if len(reqToken) < 2 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(reqToken, "Bearer ")

		if len(splitToken) < 2 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		userId, err := jwt.ParseJWT(splitToken[1], cfg.KeyToken)

		if err != nil {
			if err.Error() != jwt.InvalidToken {
				slog.Error(op+" ParseJWT", sl.Err(err))
			}

			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "id", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
