package service

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/hlog"
	"github.com/warrant-dev/warrant/pkg/config"
)

func AuthMiddleware(next http.Handler, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := hlog.FromRequest(r)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			SendErrorResponse(w, NewUnauthorizedError("Request missing Authorization header"))
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 {
			SendErrorResponse(w, NewUnauthorizedError("Invalid Authorization header"))
			logger.Warn().Msgf("Invalid Authorization header %s", authHeader)
			return
		}

		tokenType := authHeaderParts[0]
		tokenString := authHeaderParts[1]

		switch tokenType {
		case "ApiKey":
			if tokenString != config.ApiKey {
				SendErrorResponse(w, NewUnauthorizedError("Invalid API key"))
				return
			}
		default:
			SendErrorResponse(w, NewInvalidRequestError("Invalid Authorization header prefix. Must be ApiKey"))
			return
		}

		next.ServeHTTP(w, r)
	})
}