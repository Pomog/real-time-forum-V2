package service

import (
	"errors"
	"github.com/Pomog/real-time-forum-V2/internal/model"
	"github.com/Pomog/real-time-forum-V2/internal/repository"
	"github.com/Pomog/real-time-forum-V2/pkg/image"
	"time"
)

type CommentsService struct {
	repo                           repository.Comments
	postsRepo                      repository.Posts
	notificationsService           Notifications
	commentsForPage                int
	imagesDir                      string
	commentsPreModerationIsEnabled bool
}

func NewCommentsService(repo repository.Comments, postsRepo repository.Posts, notificationsService Notifications, commentsForPage int,
	imagesDir string, commentsPreModerationIsEnabled bool) *CommentsService {
	return &CommentsService{
		repo:                           repo,
		postsRepo:                      postsRepo,
		notificationsService:           notificationsService,
		commentsForPage:                commentsForPage,
		imagesDir:                      imagesDir,
		commentsPreModerationIsEnabled: commentsPreModerationIsEnabled,
	}
}

type CreateCommentInput struct {
	UserID int
	PostID int
	Data   string
	Image  string
}

func (s *CommentsService) Create(input CreateCommentInput) (model.Comment, error) {
	var comment model.Comment
	// Create comment
	imageName, err := image.Save(input.Image, s.imagesDir)
	if err != nil {
		return comment, err
	}

	comment.Author.ID = input.UserID
	comment.PostID = input.PostID
	comment.Data = input.Data
	comment.Image = imageName
	comment.Date = time.Now()

	if s.commentsPreModerationIsEnabled {
		comment.Status = model.CommentStatus.Pending
	} else {
		comment.Status = model.CommentStatus.Approved
	}

	comment.ID, err = s.repo.Create(comment)
	if err != nil {
		if errors.Is(err, repository.ErrForeignKeyConstraint) {
			return comment, ErrPostDoesntExist
		}
		return comment, err
	}

	// Create notification for post author
	post, err := s.postsRepo.GetByID(input.PostID, 0)
	if err != nil {
		if errors.Is(err, repository.ErrForeignKeyConstraint) {
			return comment, ErrPostDoesntExist
		}
		return comment, err
	}

	notification := model.Notification{
		RecipientID:  post.Author.ID,
		SenderID:     input.UserID,
		ActivityType: model.NotificationActivities.PostCommented,
		ObjectID:     input.PostID,
		Date:         time.Now(),
		Read:         false,
	}

	return comment, s.notificationsService.Create(notification)
}

func (s *CommentsService) Delete(userID, postID int) error {
	err := s.repo.Delete(userID, postID)
	if errors.Is(err, repository.ErrNoRows) {
		return ErrDeletingComment
	}

	return err
}

func (s *CommentsService) GetCommentsByPostID(postID int, userID int, page int) ([]model.Comment, error) {
	offset := (page - 1) * s.commentsForPage

	comments, err := s.repo.GetCommentsByPostID(postID, userID, s.commentsForPage, offset)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, ErrPostDoesntExist
		}
		return nil, err
	}

	return comments, nil
}

func (s *CommentsService) LikeComment(comentID, userID, likeType int) error {
	like := model.CommentLike{
		CommentID: comentID,
		UserID:    userID,
		LikeType:  likeType,
	}

	likeCreated, err := s.repo.LikeComment(like)
	if err != nil {
		if errors.Is(err, repository.ErrForeignKeyConstraint) {
			return ErrCommentDoesntExist
		}
		return err
	}

	// send notification to comment author
	if likeCreated {
		comment, err := s.repo.GetByID(comentID)
		if err != nil {
			return err
		}

		var activityType int

		if likeType == model.LikeTypes.Like {
			activityType = model.NotificationActivities.CommentLiked
		} else {
			activityType = model.NotificationActivities.CommentDisliked
		}

		notification := model.Notification{
			RecipientID:  comment.Author.ID,
			SenderID:     userID,
			ActivityType: activityType,
			Date:         time.Now(),
			Read:         false,
		}

		return s.notificationsService.Create(notification)
	}

	return nil
}
