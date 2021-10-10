package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"
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

			w.Write(models.ToBytes(http.StatusUnauthorized, "named cookie not present", nil))
			return
		}

		session, err := sm.sessionUsecase.Check(cookie.Value)
		if err != nil {
			log.Printf("error middleware 2: %v\n", err.Error())

			w.Header().Set("Content-Type", "application/json")
			w.Header().Add("Location", r.Host+"/signin") // указываем в качестве перенаправления страницу входа
			w.WriteHeader(http.StatusOK)

			w.Write(models.ToBytes(http.StatusUnauthorized, "no rights to access this resource", nil))
			return
		}

		// то есть если нашли куку и она валидна, запишем ее в контекст
		// чтобы затем использовать в последующих обработчиках
		ctx := context.WithValue(r.Context(), ContextUserId, session.UserId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "https://volchock.ru")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, Location")
		w.Header().Set("Access-Control-Max-Age", "600")
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		relativePath := r.URL.Path
		contentType := r.Header.Get("Content-Type")

		isImageUpload, _ := regexp.MatchString("^/adverts/[0-9]+/upload$", relativePath)

		switch {
		case relativePath == "/users/profile/upload", isImageUpload:
			log.Println("image upload")
			if !strings.Contains(contentType, "multipart/form-data") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(models.ToBytes(http.StatusBadRequest, "content-type: multipart/form-data required", nil))
				return
			}

		default:
			if contentType != "application/json" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(models.ToBytes(http.StatusBadRequest, "content-type: application/json required", nil))
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
