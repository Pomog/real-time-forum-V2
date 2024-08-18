package app

import (
	"github.com/Pomog/real-time-forum-V2/internal/config"
	"github.com/Pomog/real-time-forum-V2/internal/handler/http"
	"github.com/Pomog/real-time-forum-V2/internal/handler/ws"
	"github.com/Pomog/real-time-forum-V2/internal/model"
	"github.com/Pomog/real-time-forum-V2/internal/repository"
	"github.com/Pomog/real-time-forum-V2/internal/service"
	"github.com/Pomog/real-time-forum-V2/pkg/auth"
	"github.com/Pomog/real-time-forum-V2/pkg/database"
	"github.com/Pomog/real-time-forum-V2/pkg/hash"
	"github.com/Pomog/real-time-forum-V2/server"
	"log"
	"os"
)

func Run(configPath *string) {
	// Get forum config
	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Prepare database
	db, err := database.ConnectDB(
		config.Database.Driver,
		config.Database.Path,
		config.Database.FileName,
		config.Database.SchemesDir,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Run function that deletes expired sessions from database
	go repository.DeleteExpiredSessions(db)

	repos := repository.NewRepositories(db)

	passwordSalt := os.Getenv("PASSWORD_SALT")
	hasher, err := hash.NewHasher(passwordSalt)
	if err != nil {
		log.Fatalln(err)
	}

	// Prepare JWT token manager
	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	accessTokenTTL := config.AccessTokenTTL()
	refreshTokenTTL := config.RefreshTokenTTL()
	if err != nil {
		log.Fatalln(err)
	}

	tokenManager, err := auth.NewManager(jwtSigningKey, accessTokenTTL, refreshTokenTTL)
	if err != nil {
		log.Fatalln(err)
	}

	// channel to receive notifications from services and send to users by websocket
	eventsChan := make(chan *model.WSEvent)

	// Prepare services
	services := service.NewServices(service.ServicesDeps{
		Repos:                          repos,
		EventsChan:                     eventsChan,
		Hasher:                         hasher,
		TokenManager:                   tokenManager,
		AccessTokenTTL:                 accessTokenTTL,
		RefreshTokenTTL:                refreshTokenTTL,
		ImagesDir:                      config.Database.ImagesDir,
		DefaultMaleAvatar:              config.Forum.DefaultMaleAvatar,
		DefaultFemaleAvatar:            config.Forum.DefaultFemaleAvatar,
		PostsForPage:                   config.Forum.PostsForPage,
		CommentsForPage:                config.Forum.CommentsForPage,
		PostsPreModerationIsEnabled:    config.Forum.PostsPreModerationIsEnabled,
		CommentsPreModerationIsEnabled: config.Forum.CommentsPreModerationIsEnabled,
	})

	// Prepare handler
	wsHandler := ws.NewHandler(eventsChan, services, tokenManager, config)
	httpHandler := http.NewHandler(services, tokenManager, wsHandler)
	httpHandler.Init()

	// Run server
	server := server.NewServer(config, httpHandler.Router)
	log.Fatalln(server.Run())
}
