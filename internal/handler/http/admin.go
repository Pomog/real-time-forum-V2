package http

import (
	"github.com/Pomog/real-time-forum-V2/gorouter"
	"github.com/Pomog/real-time-forum-V2/internal/service"
	"net/http"
)

func (h *Handler) getRequestsForModerator(ctx *gorouter.Context) {
	requests, err := h.adminsService.GetModeratorRequests()
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	err = ctx.WriteJSON(http.StatusOK, &requests)
	if err != nil {
		return
	}
}

type RequestForModeratorActionInput struct {
	Action  string `json:"action" example:"accept"`
	Message string `json:"message"`
}

func (h *Handler) RequestForModeratorAction(ctx *gorouter.Context) {
	var input RequestForModeratorActionInput

	adminID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	requestID, err := ctx.GetIntParam("request_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err = ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	switch input.Action {
	case "accept":
		err = h.adminsService.AcceptRequestForModerator(adminID, requestID)
	case "decline":
		err = h.adminsService.DeclineRequestForModerator(adminID, requestID, input.Message)
	default:
		ctx.WriteError(http.StatusBadRequest, "invalid action")
		return
	}

	if err != nil {
		if err == service.ErrModeratorRequestDoesntExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
