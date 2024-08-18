package service

import "github.com/Pomog/real-time-forum-V2/internal/repository"

type ModeratorsService struct {
	repo repository.Moderators
}

func NewModeratorsService(repo repository.Moderators) *ModeratorsService {
	return &ModeratorsService{
		repo: repo,
	}
}
