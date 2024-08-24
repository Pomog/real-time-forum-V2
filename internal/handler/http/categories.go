package http

import (
	"errors"
	"fmt"
	"github.com/Pomog/real-time-forum-V2/gorouter"
	"github.com/Pomog/real-time-forum-V2/internal/service"
	"net/http"
)

func (h *Handler) getAllCategories(ctx *gorouter.Context) {
	categories, err := h.categoriesService.GetAll()

	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	err = ctx.WriteJSON(http.StatusOK, categories)
	fmt.Println(ctx)
	if err != nil {
		return
	}
}

func (h *Handler) getCategoryPage(ctx *gorouter.Context) {
	categoryID, err := ctx.GetIntParam("category_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	page, err := ctx.GetIntParam("page")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.categoriesService.GetByID(categoryID, page)
	if err != nil {
		if errors.Is(err, service.ErrCategoryDoesntExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = ctx.WriteJSON(http.StatusOK, category)
	if err != nil {
		return
	}
}
