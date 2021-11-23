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

	internalError "yula/internal/error"
)

var (
	logger       = logging.GetLogger()
	chatSessions = map[string]*ChatSession{} // to_string(idFrom) + "->" + to_string(idTo) => conn
)

type ChatSession struct {
	idFrom int64
	idTo   int64
	idAdv  int64

	conn []*websocket.Conn
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
	s.HandleFunc("/connect/{idFrom:[0-9]+}/{idTo:[0-9]+}/{idAdv:[0-9]+}", middleware.SetSCRFToken(http.HandlerFunc(ch.ConnectHandler))).Methods(http.MethodGet, http.MethodOptions)

	s.HandleFunc("/getDialogs/{idFrom:[0-9]+}", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(ch.getDialogsHandler)))).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/getHistory/{idFrom:[0-9]+}/{idTo:[0-9]+}/{idAdv:[0-9]+}", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(ch.getHistoryHandler)))).Methods(http.MethodGet, http.MethodOptions)

	s.Handle("/clear/{idFrom:[0-9]+}/{idTo:[0-9]+}/{idAdv:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ch.ClearHandler))).Methods(http.MethodPost, http.MethodOptions)
}

func (ch *ChatHandler) ConnectHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	vars := mux.Vars(r)
	idFrom, _ := strconv.ParseInt(vars["idFrom"], 10, 64)
	idTo, _ := strconv.ParseInt(vars["idTo"], 10, 64)
	idAdv, _ := strconv.ParseInt(vars["idAdv"], 10, 64)

	websocketConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("can not upgrade connection to websocket: ", err.Error())
		return
	}

	var curSession *ChatSession
	key := fmt.Sprintf("%d->%d:%d", idFrom, idTo, idAdv)
	if val, ok := chatSessions[key]; ok {
		curSession = val
		val.conn = append(val.conn, websocketConnection)
	} else {
		curSession = &ChatSession{
			idFrom: idFrom,
			idTo:   idTo,
			idAdv:  idAdv,

			conn: []*websocket.Conn{websocketConnection},
		}
		chatSessions[key] = curSession
	}

	go ch.HandleMessages(curSession, websocketConnection)
}

func (ch *ChatHandler) HandleMessages(session *ChatSession, conn *websocket.Conn) {
	defer func() {
		conn.Close()

		key := fmt.Sprintf("%d->%d:%d", session.idFrom, session.idTo, session.idAdv)
		for ind, value := range chatSessions[key].conn {
			if value == conn {
				chatSessions[key].conn[ind] = chatSessions[key].conn[len(chatSessions[key].conn)-1]
				chatSessions[key].conn[len(chatSessions[key].conn)-1] = nil
				chatSessions[key].conn = chatSessions[key].conn[:len(chatSessions[key].conn)-1]
			}
		}
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Debug(err)
			return
		}

		message := &models.Message{
			IdFrom: session.idFrom,
			IdTo:   session.idTo,
			IdAdv:  session.idAdv,
			Msg:    string(msg),
		}
		ch.chatUsecase.Create(message)

		key := fmt.Sprintf("%d->%d:%d", session.idTo, session.idFrom, session.idAdv)
		to := chatSessions[key]

		if to == nil {
			continue
		}

		for _, conn := range to.conn {
			if err := conn.WriteMessage(msgType, msg); err != nil {
				logger.Error("Can not write msg from user %d to user %d on ad %d", session.idFrom, session.idTo, session.idAdv)
				return
			}
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
	idAdv, _ := strconv.ParseInt(vars["idAdv"], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	messages, err := ch.chatUsecase.GetHistory(idFrom, idTo, idAdv, page.PageNum*page.Count, page.Count)
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
	idAdv, _ := strconv.ParseInt(vars["idAdv"], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = ch.chatUsecase.Clear(idFrom, idTo, idAdv)
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
