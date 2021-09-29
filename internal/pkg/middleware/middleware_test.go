package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"yula/internal/config"
	"yula/internal/models"
	sessRep "yula/internal/pkg/session/repository"
	sessUse "yula/internal/pkg/session/usecase"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_CorsMiddleware_Succsess(t *testing.T) {
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
	assert.Equal(t, string("http://89.19.190.83:5000"), w.Header().Get("Access-Control-Allow-Origin"))

}

func TestMiddleware_JsonMiddleware_Succsess(t *testing.T) {
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	caller(w, r)
	mw := JsonMiddleware(caller)
	mw.ServeHTTP(w, r)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

}

func TestMiddleware_JsonMiddleware_NoApplicationJson(t *testing.T) {
	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	caller(w, r)
	mw := JsonMiddleware(caller)
	mw.ServeHTTP(w, r)

	assert.Equal(t, string(""), w.Header().Get("Content-Type"))

}

func TestMiddleware_CheckAuthorized_Success(t *testing.T) {
	pwd, _ := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-3], "/")

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()

	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := sessUse.NewSessionUsecase(sr)

	mw := NewSessionMiddleware(su)

	userSession, err := su.Create(0)
	assert.Equal(t, true, err == nil)

	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

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
	// r := &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}

	caller(w, r)
	mw.CheckAuthorized(caller).ServeHTTP(w, r)

	cookie, err := r.Cookie("session_id")
	assert.Equal(t, nil, err)
	session, err := su.Check(cookie.Value)
	assert.Equal(t, nil, err)

	assert.Equal(t, int64(0), session.UserId)

	su.Delete(session.Value)
}

func TestMiddleware_CheckAuthorized_InvalodCookieName(t *testing.T) {
	pwd, _ := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-3], "/")

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()

	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := sessUse.NewSessionUsecase(sr)

	mw := NewSessionMiddleware(su)

	userSession, err := su.Create(0)
	assert.Equal(t, true, err == nil)

	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

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

	_, err = r.Cookie("session_id")
	assert.Equal(t, "http: named cookie not present", err.Error())
	su.Delete(userSession.Value)
}

func TestMiddleware_CheckAuthorized_InvalodCookieValue(t *testing.T) {
	pwd, _ := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-3], "/")

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()

	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := sessUse.NewSessionUsecase(sr)

	mw := NewSessionMiddleware(su)

	userSession, err := su.Create(0)
	assert.Equal(t, true, err == nil)

	caller := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

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

	_, err = r.Cookie("session_id")
	assert.Equal(t, nil, err)

	var Answer models.HttpError
	err = json.NewDecoder(w.Body).Decode(&Answer)
	if err != nil {
		t.Fatal("invalid serialization")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 401, Answer.Code)
	assert.Equal(t, "no rights to access this resource", Answer.Message)

	su.Delete(userSession.Value)
}
