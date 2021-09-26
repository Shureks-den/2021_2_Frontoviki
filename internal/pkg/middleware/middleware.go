package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"yula/internal/codes"
	"yula/internal/models"
	"yula/internal/pkg/session"

	"log"
)

type contextKey string

const ContextUserId contextKey = "user_id"

type SessionMiddleware struct {
	sessionUsecase session.SessionUsecase
}

func NewSessionMiddleware(sessionUsecase session.SessionUsecase) *SessionMiddleware {
	return &SessionMiddleware{
		sessionUsecase: sessionUsecase,
	}
}

func (sm *SessionMiddleware) CheckAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Printf("error middleware 1: %v\n", err.Error())

			w.Header().Set("Content-Type", "application/json")
			w.Header().Add("Location", r.Host+"/signin") // указываем в качестве перенаправления страницу входа
			w.WriteHeader(http.StatusOK)

			httpStat := codes.ServerErrorToHttpStatus(codes.NewServerError(codes.Unauthorized))
			response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
			js, _ := json.Marshal(response)

			w.Write(js)
			return
		}

		session, err := sm.sessionUsecase.Check(cookie.Value)
		if err != nil {
			log.Printf("error middleware 2: %v\n", err.Error())

			w.Header().Set("Content-Type", "application/json")
			w.Header().Add("Location", r.Host+"/signin") // указываем в качестве перенаправления страницу входа
			w.WriteHeader(http.StatusOK)

			httpStat := codes.ServerErrorToHttpStatus(codes.NewServerError(codes.Unauthorized))
			response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
			js, _ := json.Marshal(response)

			w.Write(js)
			return
		}

		log.Printf("session %s for user %d got", session.Value, session.UserId)

		// то есть если нашли куку и она валидна, запишем ее в контекст
		// чтобы затем использовать в последующих обработчиках
		ctx := context.WithValue(r.Context(), ContextUserId, session.UserId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
