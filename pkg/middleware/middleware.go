package middleware

import (
	"context"
	"errors"
	"filmoteka/pkg/models"
	httpResponse "filmoteka/pkg/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userId"

type Core interface {
	GetUserId(ctx context.Context, sid string) (uint64, error)
}

func AuthCheck(next http.Handler, core Core, lg *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.Response{Status: http.StatusUnauthorized, Body: nil}
		session, err := r.Cookie("session_id")

		if errors.Is(err, http.ErrNoCookie) {
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		userId, err := core.GetUserId(r.Context(), session.Value)
		if err != nil {
			lg.Errorf("auth check error: %s", err.Error())
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userId))
		if userId == 0 {
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MethodCheck(next http.Handler, method string, lg *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			response := models.Response{Status: http.StatusMethodNotAllowed, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}
		next.ServeHTTP(w, r)
	})
}
