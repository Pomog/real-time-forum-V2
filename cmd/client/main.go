package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"

	"github.com/Pomog/real-time-forum-V2/internal/config"
)

// Function for starting client server
func main() {
	var configPath = flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	indexTemp, err := template.ParseFiles("./web/public/index.html")
	if err != nil {
		log.Fatalln(err)
	}

	fileServer := http.FileServer(http.Dir("./web/src"))

	http.Handle("/src/", http.StripPrefix("/src/", fileServer))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := indexTemp.Execute(w, conf.BackendAddress())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
		}
	})

	port := conf.Client.Port
	log.Printf("Frontend server is starting at %v", port)
	if err := http.ListenAndServe(": "+port, nil); err != nil {
		log.Fatalln(err)
	}
}
