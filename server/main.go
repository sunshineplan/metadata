package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/database/mongodb/api"
	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/flags"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/sunshineplan/utils/log"
)

var (
	mongo mongodb.Client

	server = httpsvr.New()
	svc    = service.New()
)

func init() {
	svc.Name = "Metadata"
	svc.Desc = "Instance to serve Metadata"
	svc.Exec = run
	svc.TestExec = test
	svc.Options = service.Options{
		Dependencies: []string{"Wants=network-online.target", "After=network.target"},
	}
}

var (
	exclude = flag.String("exclude", "", "Exclude Files")
	logPath = flag.String("log", "", "Log Path")
)

func main() {
	self, err := os.Executable()
	if err != nil {
		svc.Fatalln("Failed to get self path:", err)
	}

	var apiClient api.Client
	flag.StringVar(&apiClient.DataSource, "source", "", "Metadata DataSource")
	flag.StringVar(&apiClient.Database, "database", "", "Metadata Database")
	flag.StringVar(&apiClient.Collection, "collection", "", "Metadata Database Collection")
	flag.StringVar(&apiClient.AppID, "id", "", "Metadata App ID")
	flag.StringVar(&apiClient.Key, "key", "", "Metadata API Key")
	flag.StringVar(&server.Unix, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "127.0.0.1", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	flag.StringVar(&svc.Options.UpdateURL, "update", "", "Update URL")
	flags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	flags.Parse()

	svc.Options.ExcludeFiles = strings.Split(*exclude, ",")
	if *logPath != "" {
		svc.SetLogger(*logPath, "", log.LstdFlags)
	}

	mongo = &apiClient
	if err := mongo.Connect(); err != nil {
		svc.Fatal(err)
	}

	if err := svc.ParseAndRun(flag.Args()); err != nil {
		svc.Fatal(err)
	}
}
