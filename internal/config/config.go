package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Pomog/real-time-forum-V2/validator"
)

type (
	Conf struct {
		API       API       `json:"api"`
		Websocket Websocket `json:"websocket"`
		Client    Client    `json:"client"`
		Database  Database  `json:"database"`
		Auth      Auth      `json:"auth"`
		Forum     Forum     `json:"forum"`
	}

	API struct {
		Host string `json:"host" validator:"required"`
		Port string `json:"port" validator:"required"`
	}

	Client struct {
		Port string `json:"port" validator:"required"`
	}

	Database struct {
		Driver     string `json:"driver" validator:"required"`
		Path       string `json:"path" validator:"required"`
		FileName   string `json:"fileName" validator:"required"`
		ImagesDir  string `json:"imagesDir" validator:"required"`
		SchemesDir string `json:"schemesDir" validator:"required"`
	}

	Auth struct {
		AccessTokenTTL  int `json:"accessTokenTTL" validator:"required"`
		RefreshTokenTTL int `json:"refreshTokenTTL" validator:"required"`
	}

	Websocket struct {
		MaxConnsForUser int   `json:"maxConnsForUser" validator:"required,max=32"`
		MaxMessageSize  int64 `json:"maxMessageSize" validator:"required"`
		TokenWait       int   `json:"tokenWait" validator:"required"`
		WriteWait       int   `json:"writeWait" validator:"required"`
		PongWait        int   `json:"pongWait" validator:"required"`
	}

	Forum struct {
		MaleAvatarsDir                 string `json:"maleAvatarsDir" validator:"required"`
		FemaleAvatarsDir               string `json:"femaleAvatarsDir" validator:"required"`
		PostsForPage                   int    `json:"postsForPage" validator:"required"`
		CommentsForPage                int    `json:"commentsForPage" validator:"required"`
		PostsPreModerationIsEnabled    bool   `json:"postsPreModerationIsEnabled"`
		CommentsPreModerationIsEnabled bool   `json:"commentsPreModerationIsEnabled"`
		DefaultFemaleAvatar            interface{}
	}
)

//Code to define all structs params

func NewConfig(confPath string) (*Conf, error) {
	var config Conf

	file, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	err = validator.Validate(config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Conf) BackendAddress() string {
	host := c.API.Host
	port := c.API.Port
	if host == "localhost" || host == "127.0.0.1" {
		return fmt.Sprintf("%s:%s", host, port)
	}
	return host
}

func (c *Conf) AccessTokenTTL() time.Duration {
	return secondsToDuration(c.Auth.AccessTokenTTL)
}

func (c *Conf) RefreshTokenTTL() time.Duration {
	return secondsToDuration(c.Auth.RefreshTokenTTL)
}

func (c *Conf) TokenWait() time.Duration {
	return secondsToDuration(c.Websocket.TokenWait)
}

func (c *Conf) WriteWait() time.Duration {
	return secondsToDuration(c.Websocket.WriteWait)
}

func (c *Conf) PongWait() time.Duration {
	return secondsToDuration(c.Websocket.PongWait)
}

func (c *Conf) PingPeriod() time.Duration {
	return (c.PongWait() * 9) / 10
}

func secondsToDuration(s int) time.Duration {
	return time.Duration(int(time.Second) * s)
}
