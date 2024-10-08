package ws

import (
	"fmt"
	"github.com/Pomog/real-time-forum-V2/internal/model"
	"sync"
)

type client struct {
	conns []*conn
	mu    sync.Mutex
	model.User
}

func (h *Handler) sendEventToClient(event *model.WSEvent) {
	client, ok := h.clients[event.RecipientID]
	if !ok {
		return
	}
	fmt.Println(event)
	client.mu.Lock()
	defer client.mu.Unlock()

	for i := 0; i < len(client.conns); i++ {
		conn := client.conns[i]

		err := conn.writeJSON(&event)
		if err != nil {
			h.closeConn(conn)
			continue
		}
	}
}

func (h *Handler) getOnlineUsers(conn *conn) error {
	var users []model.User
	for _, client := range h.clients {
		users = append(users, client.User)
	}

	return conn.writeJSON(&model.WSEvent{
		Type: model.WSEventTypes.OnlineUsersResponse,
		Body: users,
	})
}
