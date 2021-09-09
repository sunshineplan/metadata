package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/sunshineplan/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

var config mongodb.Config
var collection *mongo.Collection

func test() (err error) {
	_, err = config.Open()

	return
}

func initMongo() error {
	client, err := config.Open()
	if err != nil {
		return err
	}

	collection = client.Database(config.Database).Collection(config.Collection)

	return nil
}

func backup(file string) {
	log.Print("Start!")
	if err := initMongo(); err != nil {
		log.Fatalln("Failed to initialize mongodb:", err)
	}
	if err := config.Backup(file); err != nil {
		log.Fatal(err)
	}
	log.Print("Backup Done!")
}

func restore(file string) {
	log.Print("Start!")
	if _, err := os.Stat(file); err != nil {
		log.Fatalln("File not found:", err)
	}
	if err := config.Restore(file); err != nil {
		log.Fatal(err)
	}
	log.Print("Done!")
}

func getClientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
