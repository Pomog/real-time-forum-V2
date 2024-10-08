package service

import (
	"errors"
	"github.com/Pomog/real-time-forum-V2/internal/model"
	"github.com/Pomog/real-time-forum-V2/internal/repository"
	"time"
)

type AdminsService struct {
	repo                 repository.Admins
	notificationsService Notifications
	usersService         Users
}

func NewAdminsService(repo repository.Admins, notificationsService Notifications, usersService Users) *AdminsService {
	return &AdminsService{
		repo:                 repo,
		notificationsService: notificationsService,
		usersService:         usersService,
	}
}

func (s *AdminsService) GetModeratorRequests() ([]model.ModeratorRequest, error) {
	return s.repo.GetModeratorRequests()
}

func (s *AdminsService) AcceptRequestForModerator(adminID, requestID int) error {
	request, err := s.repo.GetModeratorRequestByID(requestID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return ErrModeratorRequestDoesntExist
		}
		return err
	}

	err = s.UpdateUserRole(request.User.ID, model.Roles.Moderator)
	if err != nil {
		return err
	}

	requestAcceptedNotification := model.Notification{
		RecipientID:  request.User.ID,
		SenderID:     adminID,
		ActivityType: model.NotificationActivities.ModeratorRequestAccepted,
		Date:         time.Now(),
		Read:         false,
	}

	err = s.notificationsService.Create(requestAcceptedNotification)
	if err != nil {
		return err
	}

	return s.DeleteModeratorRequest(requestID)
}

func (s *AdminsService) DeclineRequestForModerator(adminID, requestID int, message string) error {
	request, err := s.repo.GetModeratorRequestByID(requestID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return ErrModeratorRequestDoesntExist
		}
		return err
	}

	requestDeclinedNotification := model.Notification{
		RecipientID:  request.User.ID,
		SenderID:     adminID,
		ActivityType: model.NotificationActivities.ModeratorRequestDeclined,
		Date:         time.Now(),
		Message:      message,
		Read:         false,
	}

	err = s.notificationsService.Create(requestDeclinedNotification)
	if err != nil {
		return err
	}

	return s.DeleteModeratorRequest(requestID)
}

func (s *AdminsService) UpdateUserRole(userID int, role int) error {
	err := s.repo.UpdateUserRole(userID, role)
	if err != nil {
		return err
	}

	roleUpdatedNotification := model.Notification{
		RecipientID:  userID,
		ActivityType: model.NotificationActivities.RoleUpdated,
		Date:         time.Now(),
		Read:         false,
	}

	return s.notificationsService.Create(roleUpdatedNotification)
}

func (s *AdminsService) DeleteModeratorRequest(userID int) error {
	return s.repo.DeleteModeratorRequest(userID)
}
