package middleware

import (
	"context"
	"net/http"

	"github.com/ayeama/panel/api/internal/config"
	"github.com/ayeama/panel/api/internal/service"
)

func Session(next http.Handler, sessionService *service.SessionService, userService *service.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: because we currently apply auth to every endpoint
		ignoreURLs := []string{
			"/oauth2/callback",
		}
		for _, url := range ignoreURLs {
			if r.URL.Path == url {
				next.ServeHTTP(w, r)
				return
			}
		}

		sessionValue, err := r.Cookie(config.SessionCookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session, err := sessionService.ReadOne(sessionValue.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user := userService.ReadOne(session.UserId)

		ctx := r.Context()
		// ctx = context.WithValue(ctx, "session", session)
		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
