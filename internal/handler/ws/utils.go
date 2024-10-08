package ws

import (
	"encoding/json"
	"fmt"
	"github.com/Pomog/real-time-forum-V2/internal/model"
	"time"
)

func (h *Handler) LogConns() {
	for {
		fmt.Println(len(h.clients), "clients connected")
		for _, client := range h.clients {
			fmt.Printf("client %v have %v connections\n", client.Username, len(client.conns))
		}
		fmt.Println()
		time.Sleep(1 * time.Second)
	}
}

func unmarshalEventBody(e *model.WSEvent, v interface{}) error {
	body, ok := e.Body.(map[string]interface{})
	if !ok {
		return errInvalidEventBody
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, &v)
}
