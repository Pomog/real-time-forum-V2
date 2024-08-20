package service

import (
	"errors"
	"fmt"
	"github.com/Pomog/real-time-forum-V2/internal/model"
	"github.com/Pomog/real-time-forum-V2/internal/repository"
	"github.com/Pomog/real-time-forum-V2/pkg/auth"
	"github.com/Pomog/real-time-forum-V2/pkg/hash"
	"golang.org/x/text/language"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/cases"
)

type UsersService struct {
	repo             repository.Users
	hasher           hash.PasswordHasher
	tokenManager     auth.TokenManager
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	imagesDir        string
	maleAvatarsDir   string
	femaleAvatarsDir string
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager, accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration, imagesDir, maleAvatarsDir, femaleAvatarsDir string) *UsersService {

	return &UsersService{
		repo:             repo,
		hasher:           hasher,
		tokenManager:     tokenManager,
		accessTokenTTL:   accessTokenTTL,
		refreshTokenTTL:  refreshTokenTTL,
		imagesDir:        imagesDir,
		maleAvatarsDir:   maleAvatarsDir,
		femaleAvatarsDir: femaleAvatarsDir,
	}
}

type UsersSignUpInput struct {
	Username  string
	FirstName string
	LastName  string
	Age       int
	Gender    int
	Email     string
	Password  string
}

func (s *UsersService) SignUp(input UsersSignUpInput) error {
	var avatar string

	fmt.Println("s.maleAvatarsDir")
	fmt.Println(s.maleAvatarsDir)
	fmt.Println(s.femaleAvatarsDir)
	fmt.Println("*******************************************")

	if input.Gender == model.Gender.Male {
		avatar = s.getRandomAvatar(s.maleAvatarsDir)
	} else if input.Gender == model.Gender.Female {
		avatar = s.getRandomAvatar(s.femaleAvatarsDir)
	}

	fmt.Println("avatar")
	fmt.Println(avatar)
	fmt.Println("*******************************************")

	password := s.hasher.Hash(input.Password)

	causer := cases.Title(language.English)
	user := model.User{
		Username:   strings.ToLower(input.Username),
		FirstName:  causer.String(strings.ToLower(input.FirstName)),
		LastName:   causer.String(strings.ToLower(input.LastName)),
		Age:        input.Age,
		Gender:     input.Gender,
		Email:      strings.ToLower(input.Email),
		Password:   password,
		Registered: time.Now(),
		Role:       model.Roles.User,
		Avatar:     avatar,
	}

	err := s.repo.Create(user)

	if errors.Is(err, repository.ErrAlreadyExist) {
		return ErrUserAlreadyExist
	}

	return err
}

func (s *UsersService) getRandomAvatar(dir string) string {
	files, err := filepath.Glob(filepath.Join(dir, "*.png"))
	if err != nil || len(files) == 0 {
		return ""
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	selectedFile := files[r.Intn(len(files))]

	cleanedFilePath := strings.Replace(selectedFile, "database\\", "", -1)

	return cleanedFilePath
}

type UsersSignInInput struct {
	UsernameOrEmail string
	Password        string
}

func (s *UsersService) SignIn(input UsersSignInInput) (Tokens, error) {
	input.UsernameOrEmail = strings.ToLower(input.UsernameOrEmail)
	password := s.hasher.Hash(input.Password)

	user, err := s.repo.GetByCredentials(input.UsernameOrEmail, password)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return Tokens{}, ErrUserWrongPassword
		}
		return Tokens{}, err
	}

	return s.setSession(user.ID, user.Role)
}

func (s *UsersService) GetByID(userID int) (model.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return user, ErrUserDoesNotExist
		}
		return user, err
	}

	return user, nil
}

func (s *UsersService) GetUsersPosts(userID int) ([]model.Post, error) {
	posts, err := s.repo.GetUsersPosts(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return posts, ErrUserDoesNotExist
		}
		return posts, err
	}

	return posts, err
}

func (s *UsersService) GetUsersRatedPosts(userID int) ([]model.Post, error) {
	posts, err := s.repo.GetUsersRatedPosts(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return posts, ErrUserDoesNotExist
		}
		return posts, err
	}

	return posts, err
}

func (s *UsersService) CreateModeratorRequest(userID int) error {
	err := s.repo.CreateModeratorRequest(userID)
	if errors.Is(err, repository.ErrAlreadyExist) {
		return ErrModeratorRequestAlreadyExist
	}

	return err
}
