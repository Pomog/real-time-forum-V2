package http

import (
	"errors"
	"github.com/Pomog/real-time-forum-V2/gorouter"
	"github.com/Pomog/real-time-forum-V2/internal/service"
	"github.com/Pomog/real-time-forum-V2/pkg/image"
	"github.com/Pomog/real-time-forum-V2/validator"
	"net/http"
)

type createCommentInput struct {
	Data  string `json:"data" validator:"required,min=2,max=128"`
	Image string `json:"image"`
}

func (h *Handler) createComment(ctx *gorouter.Context) {
	var input createCommentInput

	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	postID, err := ctx.GetIntParam("post_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err = ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err = validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	newComment, err := h.commentsService.Create(service.CreateCommentInput{
		UserID: userID,
		PostID: postID,
		Data:   input.Data,
		Image:  input.Image,
	})

	if err != nil {
		if errors.Is(err, service.ErrPostDoesntExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else if errors.Is(err, image.ErrInvalidBase64String) || errors.Is(err, image.ErrTooBigImage) || errors.Is(err, image.ErrUnsupportedFormat) {
			ctx.WriteError(http.StatusBadRequest, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = ctx.WriteJSON(http.StatusCreated, newComment)
	if err != nil {
		return
	}
}

func (h *Handler) deleteComment(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	commentID, err := ctx.GetIntParam("comment_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.commentsService.Delete(userID, commentID)
	if err != nil {
		if errors.Is(err, service.ErrDeletingComment) {
			ctx.WriteError(http.StatusBadRequest, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getCommentsOfPost(ctx *gorouter.Context) {
	postID, err := ctx.GetIntParam("post_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	page, err := ctx.GetIntParam("page")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	userID, _ := ctx.GetIntParam("sub")

	comments, err := h.commentsService.GetCommentsByPostID(postID, userID, page)
	if err != nil {
		if errors.Is(err, service.ErrPostDoesntExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = ctx.WriteJSON(http.StatusOK, &comments)
	if err != nil {
		return
	}
}

type likeCommentInput struct {
	LikeType int `json:"likeType" validator:"required,min=1,max=2" example:"1"`
}

func (h *Handler) likeComment(ctx *gorouter.Context) {
	var input likeCommentInput

	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	commentID, err := ctx.GetIntParam("comment_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err = ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err = validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.commentsService.LikeComment(commentID, userID, input.LikeType)
	if err != nil {
		if errors.Is(err, service.ErrCommentDoesntExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
