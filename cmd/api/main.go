package main

import (
	"flag"
	"github.com/Pomog/real-time-forum-V2/envloader"
	"github.com/Pomog/real-time-forum-V2/internal/app"
)

func main() {
	configPath := flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	envloader.Load()

	app.Run(configPath)
}
