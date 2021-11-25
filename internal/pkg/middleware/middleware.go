package middleware

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	internalError "yula/internal/error"
	"yula/internal/models"
	session "yula/services/auth"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type contextKey string

const ContextUserId contextKey = "user_id"
const ContextLoggerField contextKey = "logger fields"

const SCRFToken = "c4e0344db55a8e7e5b79f5d2c9ff317c"

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
		ctxId := context.WithValue(r.Context(), ContextUserId, session.UserId)
		r = r.WithContext(ctxId)

		next.ServeHTTP(w, r)
	})
}

func (sm *SessionMiddleware) SoftCheckAuthorized(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, err := sm.sessionUsecase.Check(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctxId := context.WithValue(r.Context(), ContextUserId, session.UserId)
		r = r.WithContext(ctxId)
		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "https://volchock.ru")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, Location")
		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
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

		isImageUpload, _ := regexp.MatchString("^/adverts/[0-9]+/images$", relativePath)
		isImageUpload = isImageUpload && (r.Method == "POST")

		switch {
		case relativePath == "/users/profile/upload", isImageUpload:
			log.Println("image upload")
			if !strings.Contains(contentType, "multipart/form-data") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(models.ToBytes(http.StatusBadRequest, "content-type: multipart/form-data required", nil))
				return
			}

		case strings.Contains(relativePath, "/connect"):
			break

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

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		x_request_id := fmt.Sprint("", rand.Int())
		ctx := context.WithValue(r.Context(), ContextLoggerField,
			logrus.Fields{
				"x_request_id": x_request_id,
				"method":       r.Method,
				"url":          r.URL.Path,
			})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Routing(r *mux.Router) {
	r.HandleFunc("/csrf", SetSCRFToken(http.HandlerFunc(CSRFHandler))).Methods(http.MethodGet, http.MethodOptions)
}

func CSRFHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "csrf setted", nil))
}

func SetSCRFToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		next.ServeHTTP(w, r)
	}
}

func CSRFMiddleWare() func(http.Handler) http.Handler {
	space := uuid.New()
	sha := uuid.NewSHA1(space, []byte("csrf token"))
	md := uuid.NewMD5(space, []byte("csrf token"))
	token := fmt.Sprintf("%s-%s", sha.String(), md.String())

	return csrf.Protect(
		[]byte(token),
		csrf.Path("/"),
		csrf.ErrorHandler(CSRFErrorHandler()),
	)
}

func CSRFErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.CSRFErrorToken)
		w.WriteHeader(metaCode)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
	}
}
