package main

import (
	"encoding/json"
	"fmt"
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
	file := dump()
	defer os.Remove(file)
	var mailSetting mail.Setting
	c, err := query("metadata_backup")
	if err != nil {
		log.Fatal(err)
	}
	jsonbody, err := json.Marshal(c["value"])
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(jsonbody, &mailSetting); err != nil {
		log.Fatal(err)
	}
	if err := mail.SendMail(
		&mailSetting,
		fmt.Sprintf("My Metadata Backup-%s", time.Now().Format("20060102")),
		"",
		&mail.Attachment{FilePath: file, Filename: "database"},
	); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Backup My Metadata done.")
}
