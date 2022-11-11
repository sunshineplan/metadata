package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/database/mongodb/api"
	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/vharitonsky/iniflags"
)

var mongo mongodb.Client
var server = httpsvr.New()
var svc = service.Service{
	Name:     "Metadata",
	Desc:     "Instance to serve Metadata",
	Exec:     run,
	TestExec: test,
	Options:  service.Options{Dependencies: []string{"After=network.target"}},
}

var (
	exclude = flag.String("exclude", "", "Exclude Files")
	logPath = flag.String("log", "", "Log Path")
)

func main() {
	self, err := os.Executable()
	if err != nil {
		log.Fatalln("Failed to get self path:", err)
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
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.SetAllowUnknownFlags(true)
	iniflags.Parse()

	svc.Options.ExcludeFiles = strings.Split(*exclude, ",")

	mongo = &apiClient
	if err := mongo.Connect(); err != nil {
		log.Fatal(err)
	}

	if service.IsWindowsService() {
		svc.Run(false)
		return
	}

	switch flag.NArg() {
	case 0:
		run()
	case 1:
		cmd := flag.Arg(0)
		var ok bool
		if ok, err = svc.Command(cmd); !ok {
			log.Fatalln("Unknown argument:", cmd)
		}
	default:
		log.Fatalln("Unknown arguments:", strings.Join(flag.Args(), " "))
	}
	if err != nil {
		log.Fatalf("failed to %s: %v", flag.Arg(0), err)
	}
}
