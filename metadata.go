package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sunshineplan/cipher"
	"go.mongodb.org/mongo-driver/bson"
)

func query(metadata string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.FindOne(ctx, bson.M{"_id": metadata}).Decode(data)
}

func metadata(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var verify struct{ Header, Content string }
	if err := query("metadata_verify", &verify); err != nil {
		w.WriteHeader(500)
		return
	}
	var key struct{ Key string }
	if err := query("key", &key); err != nil {
		w.WriteHeader(500)
		return
	}

	header := r.Header.Get(verify.Header)
	if header == "" || header != verify.Content {
		w.WriteHeader(403)
		return
	}

	param := ps.ByName("metadata")
	if param == "key" {
		w.Write([]byte(fmt.Sprintf("%q", key.Key)))
		return
	}

	var metadata struct {
		Value     bson.M
		Allowlist []string
		Encrypt   bool
	}
	if err := query(param, &metadata); err != nil {
		w.WriteHeader(404)
		return
	}
	remote := getClientIP(r)
	if metadata.Allowlist != nil {
		var allow bool
		switch remote {
		case "127.0.0.1", "::1":
			allow = true
		case "":
			w.WriteHeader(400)
			return
		default:
			remoteIP := net.ParseIP(remote)
			for _, i := range metadata.Allowlist {
				ip, err := net.LookupIP(i)
				if err == nil {
					for _, a := range ip {
						if remoteIP.Equal(a) {
							allow = true
						}
					}
				} else {
					_, ipnet, err := net.ParseCIDR(i)
					if err != nil {
						w.WriteHeader(500)
						return
					}
					if ipnet.Contains(remoteIP) {
						allow = true
					}
				}
			}
		}
		if !allow {
			w.WriteHeader(403)
			return
		}
	}
	value, err := json.Marshal(metadata.Value)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if metadata.Encrypt {
		value = []byte(cipher.Encrypt(base64.StdEncoding.EncodeToString([]byte(key.Key)), string(value)))
	}
	w.Write(value)
	log.Printf(`- [%s] "%s" - "%s"`, remote, r.URL, r.UserAgent())
}
