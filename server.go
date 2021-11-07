package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/sunshineplan/utils/httpsvr"
)

var server httpsvr.Server

func run() {
	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	if err := initMongo(); err != nil {
		log.Fatalln("Failed to initialize mongodb:", err)
	}

	router := httprouter.New()
	router.GET("/:metadata", metadata)

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(403)
	})

	server.Handler = router
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
