package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yula/internal/models"
	"yula/internal/pkg/session"

	"log"
)

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
			log.Printf("error: %v\n", err.Error())

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := models.HttpError{Code: http.StatusUnauthorized, Message: err.Error()}
			js, _ := json.Marshal(response)

			w.Write(js)
			return
		}

		session, err := sm.sessionUsecase.Check(cookie.Value)
		if err != nil {
			log.Printf("error: %v\n", err.Error())

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := models.HttpError{Code: http.StatusUnauthorized, Message: err.Error()}
			js, _ := json.Marshal(response)

			w.Write(js)
			return
		}

		fmt.Println(cookie.Value)
		fmt.Println(session.Value, session.UserId)

		log.Printf("session for user: %d got", session.UserId)

		next.ServeHTTP(w, r)
	})
}
