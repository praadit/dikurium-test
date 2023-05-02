package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/praadit/dikurium-test/pkg/config"
)

func AuthenticationMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "claims", nil)
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		token := strings.Split(authToken, "Bearer ")
		if len(token) != 2 {
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if len(token[1]) < 1 {
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims, err := getClaims(token[1])
		if err != nil {
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx = context.WithValue(r.Context(), "claims", claims)

		handler.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}

func getClaims(accesstoken string) (*jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(accesstoken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.ApiSecret), nil
	})
	if err != nil {
		return nil, errors.New("Token ")
	}

	if !token.Valid {
		return nil, errors.New("Unauthorized")
	}

	return &claims, nil
}
