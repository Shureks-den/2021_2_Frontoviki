package delivery

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"yula/internal/models"
	"yula/internal/pkg/chat"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	internalError "yula/internal/error"
)

var (
	logger       = logging.GetLogger()
	chatSessions = map[string]*websocket.Conn{} // to_string(idFrom) + "->" + to_string(idTo) => conn
)

type ChatSession struct {
	idFrom int64
	idTo   int64

	conn *websocket.Conn
}

type ChatHandler struct {
	chatUsecase chat.ChatUsecase
}

func NewChatHandler(cu chat.ChatUsecase) *ChatHandler {
	return &ChatHandler{
		chatUsecase: cu,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ch *ChatHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	s := r.PathPrefix("/chat").Subrouter()
	s.HandleFunc("/connect/{idFrom:[0-9]+}/{idTo:[0-9]+}", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(ch.ConnectHandler)))).Methods(http.MethodGet, http.MethodOptions)

	s.HandleFunc("/getDialogs/{idFrom:[0-9]+}", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(ch.getDialogsHandler)))).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/getHistory/{idFrom:[0-9]+}/{idTo:[0-9]+}", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(ch.getHistoryHandler)))).Methods(http.MethodGet, http.MethodOptions)

	s.Handle("/clear/{idFrom:[0-9]+}/{idTo:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ch.ClearHandler))).Methods(http.MethodPost, http.MethodOptions)
}

func (ch *ChatHandler) ConnectHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	vars := mux.Vars(r)
	idFrom, _ := strconv.ParseInt(vars["idFrom"], 10, 64)
	idTo, _ := strconv.ParseInt(vars["idTo"], 10, 64)

	websocketConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("can not upgrade connection to websocket: ", err.Error())
		return
	}

	curSession := &ChatSession{
		idFrom: idFrom,
		idTo:   idTo,

		conn: websocketConnection,
	}

	key := fmt.Sprintf("%d->%d", curSession.idFrom, curSession.idTo)
	assert.Nil(nil, chatSessions[key])

	chatSessions[key] = curSession.conn

	go ch.HandleMessages(curSession)
}

func (ch *ChatHandler) HandleMessages(session *ChatSession) {
	defer func() {
		session.conn.Close()

		key := fmt.Sprintf("%d->%d", session.idFrom, session.idTo)
		delete(chatSessions, key)
	}()

	for {
		msgType, msg, err := session.conn.ReadMessage()
		if err != nil {
			logger.Debug(err)
			return
		}

		message := &models.Message{
			IdFrom: session.idFrom,
			IdTo:   session.idTo,
			Msg:    string(msg),
		}
		ch.chatUsecase.Create(message)

		key := fmt.Sprintf("%d->%d", session.idTo, session.idFrom)
		to := chatSessions[key]

		if to == nil {
			continue
		}

		if err := to.WriteMessage(msgType, msg); err != nil {
			logger.Error("Can not write msg from user %d to user %d", session.idFrom, session.idTo)
			return
		}
	}
}

func (ch *ChatHandler) getHistoryHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	query := u.Query()
	page, _ := models.NewPage(query.Get("page"), query.Get("count"))

	vars := mux.Vars(r)
	idFrom, _ := strconv.ParseInt(vars["idFrom"], 10, 64)
	idTo, err := strconv.ParseInt(vars["idTo"], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	messages, err := ch.chatUsecase.GetHistory(idFrom, idTo, page.PageNum*page.Count, page.Count)
	if err != nil {
		logger.Warnf("get history chat error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyChatHistory{Messages: messages}
	w.Write(models.ToBytes(http.StatusOK, "chat history found successfully", body))
	logger.Info("chat history found successfully")
}

func (ch *ChatHandler) ClearHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	vars := mux.Vars(r)
	idFrom, _ := strconv.ParseInt(vars["idFrom"], 10, 64)
	idTo, err := strconv.ParseInt(vars["idTo"], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = ch.chatUsecase.Clear(idFrom, idTo)
	if err != nil {
		logger.Warnf("clear chat error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "clear chat success", nil))
	logger.Info("clear chat success")
}

func (ch *ChatHandler) getDialogsHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	vars := mux.Vars(r)
	idFrom, err := strconv.ParseInt(vars["idFrom"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	dialogs, err := ch.chatUsecase.GetDialogs(idFrom)
	if err != nil {
		logger.Warnf("get dialogs error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyDialogs{Dialogs: dialogs}
	w.Write(models.ToBytes(http.StatusOK, "dialogs found successfully", body))
	logger.Info("dialogs found successfully")
}
