## To run the API server (which will create an SQLite database at ./database/forum.db on the first run), execute:
```bash
go run ./cmd/api/main.go
```

## To run the client, open a new terminal and execute:
```bash
go run ./cmd/client/main.go 
```
## Key Components of the API Server
1. Environment Loader:
   The API server loads settings from the .env file using the envloader.

2. Server Initialization:
   The Run function starts the API server using configuration from a JSON file.

3. Configuration (package config):
   The NewConfig function returns a *Conf struct that contains the configuration for various components:
```go
	Conf struct {
		API       API       `json:"api"`
		Websocket Websocket `json:"websocket"`
		Client    Client    `json:"client"`
		Database  Database  `json:"database"`
		Auth      Auth      `json:"auth"`
		Forum     Forum     `json:"forum"`
	}
```

4. Database (package database):
   Creates and returns an SQLite database connection (*sql.DB).

5. Repository Initialization (package repository):
The NewRepositories function initializes and returns a *Repositories struct that encapsulates various repositories:
```go
type Repositories struct {
	Users         Users
	Moderators    Moderators
	Admins        Admins
	Categories    Categories
	Posts         Posts
	Comments      Comments
	Notifications Notifications
	Chats         Chats
}
```

6. Hasher (package hash):
   The NewHasher function initializes a *hasher with a SALT value from the .env file:
```go
type hasher struct {
	salt string
}
```

7. JWT Manager (package auth):
   The NewManager function initializes a *Manager for handling JWT tokens, configured with values from the .env file:
```go
type Manager struct {
	jwtSigningKey   string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}
```

8. WebSocket Events (package model):
   The WSEvent struct is used for WebSocket communication:
```go
type WSEvent struct {
	Type        string      `json:"type"`
	Body        interface{} `json:"body,omitempty"`
	RecipientID int         `json:"recipientID,omitempty"`
}
```

9. Service Layer (package service):
   The NewServices function initializes a *Services struct using ServicesDeps:
```go
type Services struct {
	Users         Users
	Moderators    Moderators
	Admins        Admins
	Categories    Categories
	Posts         Posts
	Comments      Comments
	Notifications Notifications
	Chats         Chats
}

type ServicesDeps struct {
	Repos                          *repository.Repositories
	EventsChan                     chan *model.WSEvent
	Hasher                         hash.PasswordHasher
	TokenManager                   auth.TokenManager
	AccessTokenTTL                 time.Duration
	RefreshTokenTTL                time.Duration
	ImagesDir                      string
	MaleAvatarsDir                 string
	FemaleAvatarsDir               string
	PostsForPage                   int
	CommentsForPage                int
	PostsPreModerationIsEnabled    bool
	CommentsPreModerationIsEnabled bool
}
```

10. WebSocket Handlers (package ws):
    The NewHandler function initializes WebSocket handlers encapsulated in a *Handler struct:
```go
type Handler struct {
	clients         map[int]*client
	eventsChan      chan *model.WSEvent
	chatsService    service.Chats
	usersService    service.Users
	tokenManager    auth.TokenManager
	maxConnsForUser int
	maxMessageSize  int64
	tokenWait       time.Duration
	writeWait       time.Duration
	pongWait        time.Duration
	pingPeriod      time.Duration
	upgrader        websocket.Upgrader
}
```

11. HTTP Handlers (package http):
    The NewHandler function initializes HTTP handlers encapsulated in a *Handler struct:
```go
type Handler struct {
	Router               *gorouter.Router
	wsHandler            *ws.Handler
	usersService         service.Users
	moderatorsService    service.Moderators
	adminsService        service.Admins
	categoriesService    service.Categories
	postsService         service.Posts
	commentsService      service.Comments
	notificationsService service.Notifications
	chatsService         service.Chats
	tokenManager         auth.TokenManager
}
```

2.10 package server STARTS the api server
```go
type Server struct {
	httpServer *http.Server
}
```