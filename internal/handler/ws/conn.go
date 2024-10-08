package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pomog/real-time-forum-V2/internal/model"
	"github.com/Pomog/real-time-forum-V2/pkg/auth"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type conn struct {
	clientID int
	conn     *websocket.Conn
	mu       sync.Mutex
}

func (h *Handler) connReadPump(conn *conn) {
	defer h.closeConn(conn)

	conn.conn.SetReadLimit(h.maxMessageSize)
	err := conn.conn.SetReadDeadline(time.Now().Add(h.pongWait))
	if err != nil {
		return
	}

	for {
		event, err := conn.readEvent()
		if err != nil {
			err := conn.writeError(err)
			if err != nil {
				return
			}
			log.Println(err)
			return
		}

		switch event.Type {
		case model.WSEventTypes.Message:
			err = h.newMessage(conn.clientID, &event)

		case model.WSEventTypes.ChatsRequest:
			err = h.getChats(conn)

		case model.WSEventTypes.MessagesRequest:
			err = h.getMessages(conn, &event)

		case model.WSEventTypes.ReadMessageRequest:
			err = h.readMessage(conn.clientID, &event)

		case model.WSEventTypes.OnlineUsersRequst:
			err = h.getOnlineUsers(conn)

		case model.WSEventTypes.TypingInRequest:
			err = h.sendTypingInEvent(conn.clientID, &event)

		case model.WSEventTypes.PongMessage:
			err = conn.conn.SetReadDeadline(time.Now().Add(h.pongWait))

		default:
			err = errInvalidEventType
		}

		if err != nil {
			log.Println(err.Error())
			err := conn.writeError(err)
			if err != nil {
				return
			}
			return
		}
	}
}

func (h *Handler) pingConn(c *conn) {
	ticker := time.NewTicker(h.pingPeriod)
	defer func() {
		ticker.Stop()
		h.closeConn(c)
	}()
	for {
		<-ticker.C
		err := c.conn.SetWriteDeadline(time.Now().Add(h.pongWait))
		if err != nil {
			return
		}

		event := &model.WSEvent{
			Type: model.WSEventTypes.PingMessage,
		}

		err = c.writeJSON(event)
		if err != nil {
			return
		}
	}
}

func (h *Handler) closeConn(c *conn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, ok := h.clients[c.clientID]
	if !ok {
		return
	}

	for i := 0; i < len(client.conns); i++ {
		if client.conns[i] == c {
			err := client.conns[i].conn.Close()
			if err != nil {
				return
			}
			client.conns = append(client.conns[:i], client.conns[i+1:]...)
			break
		}
	}

	if len(client.conns) == 0 {
		delete(h.clients, client.ID)
	}
}

func (c *conn) writeJSON(data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteJSON(data)
}

func (c *conn) writeError(err error) error {
	return c.writeJSON(
		&model.WSEvent{
			Type: model.WSEventTypes.Error,
			Body: err.Error(),
		},
	)
}

func (c *conn) readEvent() (model.WSEvent, error) {
	var event model.WSEvent

	_, messageBytes, err := c.conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		return event, err
	}

	err = json.Unmarshal(messageBytes, &event)

	return event, err
}

func (h *Handler) identifyConn(c *conn) error {
	err := c.conn.SetReadDeadline(time.Now().Add(h.tokenWait))
	if err != nil {
		return err
	}

	_, messageBytes, err := c.conn.ReadMessage()
	if err != nil {
		return errNoTokenReceived
	}

	var event model.WSEvent
	err = json.Unmarshal(messageBytes, &event)
	if err != nil {
		return err
	}

	token := fmt.Sprintf("%s", event.Body)
	sub, _, err := h.tokenManager.Parse(token)
	if err != nil {
		if errors.Is(err, auth.ErrExpiredToken) {
			err := c.writeError(err)
			if err != nil {
				return err
			}
			return h.identifyConn(c)
		}
		return err
	}

	c.clientID = sub
	return nil
}
