package http

import (
	"github.com/Pomog/real-time-forum-V2/gorouter"
	"net/http"
)

func (h *Handler) getNotifications(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	notifications, err := h.notificationsService.GetNotifications(userID)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
	}

	err = ctx.WriteJSON(http.StatusOK, notifications)
	if err != nil {
		return
	}
}
