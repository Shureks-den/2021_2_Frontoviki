package chat

import (
	"net/http"
	"strconv"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	logger   logging.Logger = logging.GetLogger()
	sessions []*ChatSession
)

type ChatSession struct {
	idFrom int64
	idTo   int64

	conn *websocket.Conn
}

type ChatHandler struct {
	// chatUsecase chat.ChatUsecase
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		// chatUsecase: chat.ChatUsecase,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ch *ChatHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	s := r.PathPrefix("/chat").Subrouter()
	// s.HandleFunc("/getHistory/{idFrom:[0-9]+}/{idTo:[0-9]+}", ch.getHistory).Methods(http.MethodGet, http.MethodOptions)

	s.Handle("/connect/{idFrom:[0-9]+}/{idTo:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ch.Connect))).Methods(http.MethodGet, http.MethodOptions)
}

func (ch *ChatHandler) Connect(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	vars := mux.Vars(r)
	idFrom, _ := strconv.ParseInt(vars["idFrom"], 10, 64)

	peer, _ := upgrader.Upgrade(w, r, nil)

	chatSession := ChatSession{
		idFrom: idFrom,
		idTo:   0,
		conn:   peer,
	}

	sessions = append(sessions, &chatSession)

	go func() {
		for {
			msgType, msg, err := chatSession.conn.ReadMessage()
			if err != nil {
				logger.Error("Can not read msg from user %d to user %d", chatSession.idFrom, chatSession.idTo)
				return
			}

			if err := chatSession.conn.WriteMessage(msgType, msg); err != nil {
				logger.Error("Can not write msg from user %d to user %d", chatSession.idFrom, chatSession.idTo)
				return
			}
		}
	}()
}

// func (ch *ChatHandler) getHistory(w http.ResponseWriter, r *http.Request) {
// 	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

// 	u, err := url.Parse(r.URL.RequestURI())
// 	if err != nil {
// 		w.WriteHeader(http.StatusOK)
// 		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
// 		w.Write(models.ToBytes(metaCode, metaMessage, nil))
// 		return
// 	}

// 	query := u.Query()
// 	page, err := models.NewPage(query.Get("page"), query.Get("count"))

// 	vars := mux.Vars(r)
// 	idFrom, err := strconv.ParseInt(vars["idFrom"], 10, 64)
// 	idTo, err := strconv.ParseInt(vars["idTo"], 10, 64)

// 	if err != nil {
// 		w.WriteHeader(http.StatusOK)
// 		metaCode, metaMessage := internalError.ToMetaStatus(err)
// 		w.Write(models.ToBytes(metaCode, metaMessage, nil))
// 		return
// 	}

// 	messages, err := ch.chatUsecase.GetHistory(idFrom, idTo, page.PageNum, page.Count)
// 	if err != nil {
// 		logger.Warnf("get history chat error: %s", err.Error())
// 		w.WriteHeader(http.StatusOK)

// 		metaCode, metaMessage := internalError.ToMetaStatus(err)
// 		w.Write(models.ToBytes(metaCode, metaMessage, nil))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	body := models.HttpBodyChatHistory{Messages: messages}
// 	w.Write(models.ToBytes(http.StatusOK, "chat history found successfully", body))
// 	logger.Info("chat history found successfully")
// }
