package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sunshineplan/utils/ste"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func metadata(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	verify, err := query("metadata_verify")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	header := r.Header.Get(verify["header"].(string))
	if header == "" || header != verify["content"] {
		w.WriteHeader(403)
		return
	}
	metadata, err := query(ps.ByName("metadata"))
	if err != nil {
		w.WriteHeader(404)
		return
	}
	allowlist := metadata["allowlist"]
	remote := getClientIP(r)
	if allowlist != nil {
		var allow bool
		switch remote {
		case "127.0.0.1", "::1":
			allow = true
		case "":
			w.WriteHeader(400)
			return
		default:
			remoteIP := net.ParseIP(remote)
			for _, i := range allowlist.(primitive.A) {
				ip, err := net.LookupIP(i.(string))
				if err == nil {
					for _, a := range ip {
						if remoteIP.Equal(a) {
							allow = true
						}
					}
				} else {
					_, ipnet, err := net.ParseCIDR(i.(string))
					if err != nil {
						w.WriteHeader(500)
						return
					}
					if ipnet.Contains(remoteIP) {
						allow = true
					}
				}
			}
			if !allow {
				w.WriteHeader(403)
				return
			}
		}
	}
	value, err := json.Marshal(metadata["value"])
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if metadata["encrypt"] == 1 {
		key, err := query("key")
		if err != nil || key["value"] == nil {
			w.WriteHeader(500)
			return
		}
		w.Write([]byte(ste.Encrypt(base64.StdEncoding.EncodeToString([]byte(key["value"].(string))), string(value))))
		log.Printf(`- [%s] "%s" - "%s"`, remote, r.URL, r.UserAgent())
		return
	}
	w.Write(value)
	log.Printf(`- [%s] "%s" - "%s"`, remote, r.URL, r.UserAgent())
}
