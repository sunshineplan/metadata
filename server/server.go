package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func test() error {
	var key struct{ Key string }
	return query("key", &key)
}

func run() error {
	router := httprouter.New()
	router.GET("/:metadata", metadata)

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(403)
	})

	server.Handler = router
	return server.Run()
}
