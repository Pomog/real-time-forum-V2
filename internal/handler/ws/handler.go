package ws

import (
	"github.com/Pomog/real-time-forum-V2/gorouter"
	"github.com/Pomog/real-time-forum-V2/internal/config"
	"github.com/Pomog/real-time-forum-V2/internal/model"
	"github.com/Pomog/real-time-forum-V2/internal/service"
	"github.com/Pomog/real-time-forum-V2/pkg/auth"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Handler struct {
	clients         map[int]*client
	eventsChan      chan *model.WSEvent
	chatsService    service.Chats
	usersService    service.Users
	tokenManager    auth.TokenManager
	maxConnsForUser int
	maxMessageSize  int64
	tokenWait       time.Duration
	writeWait       time.Duration
	pongWait        time.Duration
	pingPeriod      time.Duration
	upgrader        websocket.Upgrader
}

func NewHandler(eventsChan chan *model.WSEvent, services *service.Services, tokenManager auth.TokenManager, config *config.Conf) *Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return &Handler{
		clients:         make(map[int]*client),
		eventsChan:      eventsChan,
		chatsService:    services.Chats,
		usersService:    services.Users,
		tokenManager:    tokenManager,
		upgrader:        upgrader,
		maxConnsForUser: config.Websocket.MaxConnsForUser,
		maxMessageSize:  config.Websocket.MaxMessageSize,
		tokenWait:       config.TokenWait(),
		writeWait:       config.WriteWait(),
		pongWait:        config.PongWait(),
		pingPeriod:      config.PingPeriod(),
	}
}

func (h *Handler) ServeWS(ctx *gorouter.Context) {
	ws, err := h.upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	connection := &conn{conn: ws}

	err = h.identifyConn(connection)
	if err != nil {
		err := connection.writeError(err)
		if err != nil {
			return
		}
		err = connection.conn.Close()
		if err != nil {
			return
		}
		return
	}

	c, ok := h.clients[connection.clientID]
	if !ok {
		user, err := h.usersService.GetByID(connection.clientID)
		if err != nil {
			err := connection.writeError(err)
			if err != nil {
				return
			}
			err = connection.conn.Close()
			if err != nil {
				return
			}
			return
		}

		c = &client{User: user}
		h.clients[connection.clientID] = c
	}

	if len(c.conns) == h.maxConnsForUser {
		conn := c.conns[0]
		err = conn.writeError(errTooManyConnections)
		if err != nil {
			return
		}
		h.closeConn(conn)
	}

	err = connection.writeJSON(&model.WSEvent{Type: model.WSEventTypes.SuccessConnection})
	if err != nil {
		return
	}

	go h.connReadPump(connection)
	go h.pingConn(connection)

	c.conns = append(c.conns, connection)
}

func (h *Handler) RunEventsPump() {
	for {
		event, ok := <-h.eventsChan
		if !ok {
			continue
		}
		h.sendEventToClient(event)
	}
}
