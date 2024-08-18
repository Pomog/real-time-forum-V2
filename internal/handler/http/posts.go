package http

import (
	"errors"
	"github.com/Pomog/real-time-forum-V2/gorouter"
	"github.com/Pomog/real-time-forum-V2/internal/service"
	"github.com/Pomog/real-time-forum-V2/pkg/image"
	"github.com/Pomog/real-time-forum-V2/validator"
	"net/http"
)

func (h *Handler) getPost(ctx *gorouter.Context) {
	postID, err := ctx.GetIntParam("post_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	userID, _ := ctx.GetIntParam("sub")

	post, err := h.postsService.GetByID(postID, userID)
	if err != nil {
		if errors.Is(err, service.ErrPostDoesntExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = ctx.WriteJSON(http.StatusOK, &post)
	if err != nil {
		return
	}
}

type createPostInput struct {
	Title      string `json:"title" validator:"required,min=2,max=64"`
	Data       string `json:"data" validator:"required,min=2,max=512"`
	Image      string `json:"image"`
	Categories []int  `json:"categories" validator:"required,min=1" example:"1,2"`
}

type createPostResponse struct {
	PostID int `json:"postID" example:"1"`
}

func (h *Handler) createPost(ctx *gorouter.Context) {
	var input createPostInput
	userID, err := ctx.GetIntParam("sub")
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

	newPostID, err := h.postsService.Create(service.CreatePostInput{
		UserID:     userID,
		Title:      input.Title,
		Data:       input.Data,
		Image:      input.Image,
		Categories: input.Categories,
	})

	if err != nil {
		if errors.Is(err, service.ErrCategoryDoesntExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else if errors.Is(err, service.ErrTooManyCategories) {
			ctx.WriteError(http.StatusBadRequest, err.Error())
		} else if errors.Is(err, image.ErrInvalidBase64String) || errors.Is(err, image.ErrTooBigImage) || errors.Is(err, image.ErrUnsupportedFormat) {
			ctx.WriteError(http.StatusBadRequest, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := createPostResponse{PostID: newPostID}
	err = ctx.WriteJSON(http.StatusCreated, &resp)
	if err != nil {
		return
	}
}

func (h *Handler) deletePost(ctx *gorouter.Context) {
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

	if err = h.postsService.Delete(userID, postID); err != nil {
		if errors.Is(err, service.ErrDeletingPost) {
			ctx.WriteError(http.StatusBadRequest, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusNoContent)
}

type likePostInput struct {
	LikeType int `json:"likeType" validator:"required,min=1,max=2" example:"1"`
}

func (h *Handler) likePost(ctx *gorouter.Context) {
	var input likePostInput

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

	err = h.postsService.LikePost(postID, userID, input.LikeType)
	if err != nil {
		if errors.Is(err, service.ErrPostDoesntExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
