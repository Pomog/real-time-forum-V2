package service

import (
	"errors"
	"github.com/Pomog/real-time-forum-V2/internal/model"
	"github.com/Pomog/real-time-forum-V2/internal/repository"
	"github.com/Pomog/real-time-forum-V2/pkg/auth"
	"time"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UsersRefreshTokensInput struct {
	AccessToken  string
	RefreshToken string
}

func (s *UsersService) RefreshTokens(input UsersRefreshTokensInput) (Tokens, error) {
	sub, role, err := s.tokenManager.Parse(input.AccessToken)
	if err != nil {
		if !errors.Is(err, auth.ErrExpiredToken) {
			return Tokens{}, err
		}
	}

	err = s.repo.DeleteSession(sub, input.RefreshToken)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return Tokens{}, ErrSessionNotFound
		}
		return Tokens{}, err
	}

	return s.setSession(sub, role)
}

func (s *UsersService) setSession(userID, role int) (Tokens, error) {
	accessToken, err := s.tokenManager.NewJWT(userID, role)
	refreshToken := s.tokenManager.NewRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	session := model.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return tokens, s.repo.SetSession(session)
}
