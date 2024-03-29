package main

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sunshineplan/cipher"
	"github.com/sunshineplan/database/mongodb"
)

func query(metadata string, data any) error {
	return mongo.FindOne(mongodb.M{"_id": metadata}, nil, data)
}

func metadata(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var verify struct{ Header, Content string }
	if err := query("metadata_verify", &verify); err != nil {
		svc.Print(err)
		w.WriteHeader(500)
		return
	}
	if header := r.Header.Get(verify.Header); header != verify.Content {
		w.WriteHeader(403)
		return
	}

	param := ps.ByName("metadata")
	var metadata struct {
		Value     any
		Allowlist []string
		Encrypt   bool
	}
	if err := query(param, &metadata); err != nil {
		svc.Print(err)
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
						svc.Print(err)
						continue
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
	value, err := json.Marshal(&metadata.Value)
	if err != nil {
		svc.Print(err)
		w.WriteHeader(500)
		return
	}
	if metadata.Encrypt && param != "key" {
		var key struct{ Value string }
		if err := query("key", &key); err != nil {
			svc.Print(err)
			w.WriteHeader(500)
			return
		}

		value = []byte(cipher.EncryptText(base64.StdEncoding.EncodeToString([]byte(key.Value)), string(value)))
	}
	w.Write(value)
	svc.Printf(`- [%s] "%s" - "%s"`, remote, r.URL, r.UserAgent())
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
