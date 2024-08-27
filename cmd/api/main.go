package main

import (
	"flag"

	"github.com/Pomog/real-time-forum-V2/envloader"
	"github.com/Pomog/real-time-forum-V2/internal/app"
)

// Function for geting commands from terminal to start API server
func main() {
	configPath := flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	envloader.Load()

	app.Run(configPath)
}
