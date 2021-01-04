package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sunshineplan/utils/mail"
)

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

func backup() {
	log.Print("Start!")
	tmpfile, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	tmpfile.Close()
	if err := mongo.Backup(tmpfile.Name()); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	var backup struct {
		Value struct {
			From, SMTPServer, Password string
			SMTPServerPort             int
			To                         []string
		}
	}
	if err := query("metadata_backup", &backup); err != nil {
		log.Fatal(err)
	}

	if err := (&mail.Dialer{
		Host:     backup.Value.SMTPServer,
		Port:     backup.Value.SMTPServerPort,
		Account:  backup.Value.From,
		Password: backup.Value.Password,
	}).Send(&mail.Message{
		To:          backup.Value.To,
		Subject:     fmt.Sprintf("Metadata Backup-%s", time.Now().Format("20060102")),
		Attachments: []*mail.Attachment{{Path: tmpfile.Name(), Filename: "database"}},
	}); err != nil {
		log.Fatal(err)
	}
	log.Print("Backup Done!")
}
