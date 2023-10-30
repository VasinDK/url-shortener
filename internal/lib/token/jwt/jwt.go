package jwt

import (
	"fmt"
	"mod_shortener/internal/lib/api/user"
	"mod_shortener/internal/lib/logger/sl"
	"time"

	"log/slog"

	"github.com/golang-jwt/jwt/v4"
)

const InvalidToken = "Invalid token"

type CastomClaim struct {
	*jwt.RegisteredClaims
	Name  string
	Login string
	Email string
}

func CreateJWT(AccessTTL int64, KeyToken string, user *user.User) (string, error) {
	const op = "JWT.CreateJWT"

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, CastomClaim{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(AccessTTL))),
			ID:        user.Id,
		},
		user.Name,
		user.Login,
		user.Email,
	})

	token, err := t.SignedString([]byte(KeyToken))

	if err != nil {
		slog.Error(op+".SignedString", sl.Err(err))
		return "", err
	}

	return token, nil
}

func CreateRefreshToken() {}

func ParseJWT(jwtToken, KeyToken string) (string, error) {
	var userClaim CastomClaim

	token, err := jwt.ParseWithClaims(jwtToken, &userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(KeyToken), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf(InvalidToken)
	}

	return userClaim.ID, nil
}
