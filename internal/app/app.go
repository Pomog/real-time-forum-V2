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
	configServer, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.ConnectDB(
		configServer.Database.Driver,
		configServer.Database.Path,
		configServer.Database.FileName,
		configServer.Database.SchemesDir,
	)
	if err != nil {
		log.Fatalln(err)
	}

	go repository.DeleteExpiredSessions(db)

	repos := repository.NewRepositories(db)

	passwordSalt := os.Getenv("PASSWORD_SALT")
	hasher, err := hash.NewHasher(passwordSalt)
	if err != nil {
		log.Fatalln(err)
	}

	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	accessTokenTTL := configServer.AccessTokenTTL()
	refreshTokenTTL := configServer.RefreshTokenTTL()
	if err != nil {
		log.Fatalln(err)
	}

	tokenManager, err := auth.NewManager(jwtSigningKey, accessTokenTTL, refreshTokenTTL)
	if err != nil {
		log.Fatalln(err)
	}

	eventsChan := make(chan *model.WSEvent)

	services := service.NewServices(service.ServicesDeps{
		Repos:                          repos,
		EventsChan:                     eventsChan,
		Hasher:                         hasher,
		TokenManager:                   tokenManager,
		AccessTokenTTL:                 accessTokenTTL,
		RefreshTokenTTL:                refreshTokenTTL,
		ImagesDir:                      configServer.Database.ImagesDir,
		MaleAvatarsDir:                 configServer.Forum.MaleAvatarsDir,
		FemaleAvatarsDir:               configServer.Forum.FemaleAvatarsDir,
		PostsForPage:                   configServer.Forum.PostsForPage,
		CommentsForPage:                configServer.Forum.CommentsForPage,
		PostsPreModerationIsEnabled:    configServer.Forum.PostsPreModerationIsEnabled,
		CommentsPreModerationIsEnabled: configServer.Forum.CommentsPreModerationIsEnabled,
	})

	wsHandler := ws.NewHandler(eventsChan, services, tokenManager, configServer)
	httpHandler := http.NewHandler(services, tokenManager, wsHandler)
	httpHandler.Init()

	serverAPI := server.NewServer(configServer, httpHandler.Router)
	log.Fatalln(serverAPI.Run())
}
