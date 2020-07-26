package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
)

func run() {
	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	router := httprouter.New()
	router.GET("/:metadata", metadata)
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(403)
	})

	if *unix != "" && OS == "linux" {
		if _, err := os.Stat(*unix); err == nil {
			err = os.Remove(*unix)
			if err != nil {
				log.Fatalf("Failed to remove socket file: %v", err)
			}
		}

		listener, err := net.Listen("unix", *unix)
		if err != nil {
			log.Fatalf("Failed to listen socket file: %v", err)
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			if err := listener.Close(); err != nil {
				log.Printf("Failed to close listener: %v", err)
			}
			if _, err := os.Stat(*unix); err == nil {
				if err := os.Remove(*unix); err != nil {
					log.Printf("Failed to remove socket file: %v", err)
				}
			}
			close(idleConnsClosed)
		}()

		if err := os.Chmod(*unix, 0666); err != nil {
			log.Fatalf("Failed to chmod socket file: %v", err)
		}

		http.Serve(listener, router)
		<-idleConnsClosed
	} else {
		http.ListenAndServe(*host+":"+*port, router)
	}
}
