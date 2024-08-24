package http

import (
	"errors"
	"fmt"
	"github.com/Pomog/real-time-forum-V2/gorouter"
	"github.com/Pomog/real-time-forum-V2/internal/service"
	"github.com/Pomog/real-time-forum-V2/validator"
	"net/http"
)

type usersSignUpInput struct {
	Username  string `json:"username" validator:"required,username,min=2,max=64" example:"johndoe"`
	FirstName string `json:"firstName" validator:"required,min=2,max=64" example:"John"`
	LastName  string `json:"lastName" validator:"required,min=2,max=64" example:"Doe"`
	Age       int    `json:"age" validator:"required,min=12,max=110" example:"18"`
	Gender    int    `json:"gender" validator:"min=1,max=2" example:"1"`
	Email     string `json:"email" validator:"required,email,max=64" example:"johndoe@gmail.com"`
	Password  string `json:"password" validator:"required,password,min=7,max=64" example:"Password123@"`
}

func (h *Handler) usersSignUp(ctx *gorouter.Context) {
	var input usersSignUpInput

	fmt.Println("IN the func (h *Handler) usersSignUp(ctx *gorouter.Context) {")

	if err := ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(input)

	if err := validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err := h.usersService.SignUp(service.UsersSignUpInput{
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Age:       input.Age,
		Gender:    input.Gender,
		Email:     input.Email,
		Password:  input.Password,
	})

	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExist) {
			ctx.WriteError(http.StatusConflict, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusCreated)
}

type usersSignInInput struct {
	UsernameOrEmail string `json:"usernameOrEmail" validator:"required,max=64" example:"johndoe"`
	Password        string `json:"password" validator:"required,max=64" example:"Password123@"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) usersSignIn(ctx *gorouter.Context) {
	var input usersSignInInput

	fmt.Println("IN the func (h *Handler) usersSignIn")

	if err := ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.usersService.SignIn(service.UsersSignInInput{
		UsernameOrEmail: input.UsernameOrEmail,
		Password:        input.Password,
	})

	if err != nil {
		if errors.Is(err, service.ErrUserWrongPassword) {
			ctx.WriteError(http.StatusUnauthorized, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	err = ctx.WriteJSON(http.StatusOK, resp)
	if err != nil {
		return
	}
}

func (h *Handler) getUser(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("user_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.usersService.GetByID(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = ctx.WriteJSON(http.StatusOK, user)
	if err != nil {
		return
	}
}

func (h *Handler) getUsersPosts(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("user_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	posts, err := h.usersService.GetUsersPosts(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = ctx.WriteJSON(http.StatusOK, posts)
	if err != nil {
		return
	}
}

func (h *Handler) getUsersRatedPosts(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("user_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	posts, err := h.usersService.GetUsersRatedPosts(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExist) {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = ctx.WriteJSON(http.StatusOK, posts)
	if err != nil {
		return
	}
}

type usersRefreshTokensInput struct {
	AccessToken  string `json:"accessToken" validator:"required"`
	RefreshToken string `json:"refreshToken" validator:"required"`
}

func (h *Handler) usersRefreshTokens(ctx *gorouter.Context) {
	var input usersRefreshTokensInput

	if err := ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.usersService.RefreshTokens(service.UsersRefreshTokensInput{
		AccessToken:  input.AccessToken,
		RefreshToken: input.RefreshToken,
	})

	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			ctx.WriteError(http.StatusUnauthorized, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	err = ctx.WriteJSON(http.StatusOK, resp)
	if err != nil {
		return
	}
}

func (h *Handler) requestModerator(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.usersService.CreateModeratorRequest(userID)
	if err != nil {
		if errors.Is(err, service.ErrModeratorRequestAlreadyExist) {
			ctx.WriteError(http.StatusConflict, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusCreated)
}
