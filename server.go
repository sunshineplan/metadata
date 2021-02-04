package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sunshineplan/cipher"
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

	router := httprouter.New()
	router.GET("/:metadata", metadata)
	router.POST("/do", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		mode := r.FormValue("mode")
		key := r.FormValue("key")
		content := r.FormValue("content")
		switch mode {
		case "encrypt":
			w.Header().Set("Content-Type", "application/json")
			data, _ := json.Marshal(map[string]interface{}{"result": cipher.Encrypt(key, content)})
			w.Write(data)
		case "decrypt":
			w.Header().Set("Content-Type", "application/json")
			result, err := cipher.Decrypt(key, strings.TrimSpace(content))
			var data []byte
			if err != nil {
				data, _ = json.Marshal(map[string]interface{}{"result": nil})
			} else {
				data, _ = json.Marshal(map[string]interface{}{"result": result})
			}
			w.Write(data)
		default:
			w.WriteHeader(400)
		}
	})

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(403)
	})

	server.Handler = router
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
