package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"yula/internal/models"

	myerr "yula/internal/error"

	sessMock "yula/services/auth/mocks"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_CorsMiddleware_Success(t *testing.T) {
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	caller(w, r)
	mw := CorsMiddleware(caller)
	mw.ServeHTTP(w, r)

	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, string("Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, Location"),
		w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, string("POST, GET, OPTIONS, PUT, DELETE"), w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, string("https://volchock.ru"), w.Header().Get("Access-Control-Allow-Origin"))

}

func TestMiddleware_JsonMiddleware_Success(t *testing.T) {
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	caller(w, r)
	mw := ContentTypeMiddleware(caller)
	mw.ServeHTTP(w, r)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

}

func TestMiddleware_JsonMiddleware_NoApplicationJson(t *testing.T) {
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	caller(w, r)
	mw := ContentTypeMiddleware(caller)
	mw.ServeHTTP(w, r)

	assert.Equal(t, string("application/json"), w.Header().Get("Content-Type"))

}

func TestMiddleware_CheckAuthorized_Success(t *testing.T) {
	su := sessMock.SessionUsecase{}
	mw := NewSessionMiddleware(&su)
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	userSession := models.Session{Value: uuid.NewString(), UserId: 0, ExpiresAt: time.Now().Add(time.Hour)}

	su.On("Check", userSession.Value).Return(&userSession, nil)

	w := httptest.NewRecorder()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userSession.Value,
		Expires:  userSession.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})
	r := &http.Request{Header: http.Header{"Cookie": []string{w.Header().Get("Set-Cookie")}}}

	caller(w, r)
	mw.CheckAuthorized(caller).ServeHTTP(w, r)

	cookie, err := r.Cookie("session_id")
	assert.Nil(t, err)
	session, err := su.Check(cookie.Value)
	assert.Nil(t, err)

	assert.Equal(t, int64(0), session.UserId)
}

func TestMiddleware_CheckSoftAuthorized_Success(t *testing.T) {
	su := sessMock.SessionUsecase{}
	mw := NewSessionMiddleware(&su)
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	userSession := models.Session{Value: uuid.NewString(), UserId: 0, ExpiresAt: time.Now().Add(time.Hour)}

	su.On("Check", userSession.Value).Return(&userSession, nil)

	w := httptest.NewRecorder()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userSession.Value,
		Expires:  userSession.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})
	r := &http.Request{Header: http.Header{"Cookie": []string{w.Header().Get("Set-Cookie")}}}

	caller(w, r)
	mw.SoftCheckAuthorized(caller).ServeHTTP(w, r)

	cookie, err := r.Cookie("session_id")
	assert.Nil(t, err)
	session, err := su.Check(cookie.Value)
	assert.Nil(t, err)

	assert.Equal(t, int64(0), session.UserId)
}

func TestMiddleware_CheckAuthorized_InvalidCookieName(t *testing.T) {
	su := sessMock.SessionUsecase{}
	mw := NewSessionMiddleware(&su)
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	userSession := models.Session{Value: uuid.NewString(), UserId: 0, ExpiresAt: time.Now().Add(time.Hour)}

	su.On("Check", userSession.Value).Return(&userSession, nil)

	w := httptest.NewRecorder()
	http.SetCookie(w, &http.Cookie{
		Name:     "wrong_session_id",
		Value:    userSession.Value,
		Expires:  userSession.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})
	r := &http.Request{Header: http.Header{"Cookie": []string{w.Header().Get("Set-Cookie")}}}

	caller(w, r)
	mw.CheckAuthorized(caller).ServeHTTP(w, r)

	_, err := r.Cookie("session_id")
	assert.Equal(t, "http: named cookie not present", err.Error())
}

func TestMiddleware_CheckAuthorized_InvalidCookieValue(t *testing.T) {
	su := sessMock.SessionUsecase{}
	mw := NewSessionMiddleware(&su)
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	userSession := models.Session{Value: uuid.NewString(), UserId: 0, ExpiresAt: time.Now().Add(time.Hour)}

	su.On("Check", "wrong_value").Return(nil, myerr.NotExist)

	w := httptest.NewRecorder()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "wrong_value",
		Expires:  userSession.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})
	r := &http.Request{Header: http.Header{"Cookie": []string{w.Header().Get("Set-Cookie")}}}

	caller(w, r)
	mw.CheckAuthorized(caller).ServeHTTP(w, r)

	_, err := r.Cookie("session_id")
	assert.Nil(t, err)

	var Answer models.HttpError
	err = json.NewDecoder(w.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 401, Answer.Code)
	assert.Equal(t, "no rights to access this resource", Answer.Message)

}

func TestLoggerInit(t *testing.T) {
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mw := LoggerMiddleware(caller)
	mw.ServeHTTP(w, r)
}

func TestSetCSRF(t *testing.T) {
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mw := SetSCRFToken(caller)
	mw.ServeHTTP(w, r)
}

func TestCSRF(t *testing.T) {
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mw := CSRFMiddleWare()
	router := mux.NewRouter()
	router.Use(mw)
	caller(w, r)
}
