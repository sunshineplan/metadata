package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunshineplan/utils/service"
	"github.com/vharitonsky/iniflags"
)

var logPath *string

var svc = service.Service{
	Name:    "Metadata",
	Desc:    "Instance to serve Metadata",
	Exec:    run,
	Options: service.Options{Dependencies: []string{"After=network.target"}},
}

func main() {
	self, err := os.Executable()
	if err != nil {
		log.Fatalln("Failed to get self path:", err)
	}

	flag.StringVar(&mongo.Server, "dbserver", "localhost", "Metadata Database Server Address")
	flag.IntVar(&mongo.Port, "dbport", 27017, "Metadata Database Port")
	flag.StringVar(&mongo.Database, "database", "", "Metadata Database Name")
	flag.StringVar(&mongo.Collection, "collection", "", "Metadata Database Collection Name")
	flag.StringVar(&mongo.Username, "username", "", "Metadata Database Username")
	flag.StringVar(&mongo.Password, "password", "", "Metadata Database Password")
	flag.StringVar(&server.Unix, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "127.0.0.1", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	logPath = flag.String("log", "", "Log Path")
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	if service.IsWindowsService() {
		svc.Run(false)
		return
	}

	switch flag.NArg() {
	case 0:
		run()
	case 1:
		switch flag.Arg(0) {
		case "run", "debug":
			run()
		case "install":
			err = svc.Install()
		case "remove":
			err = svc.Remove()
		case "start":
			err = svc.Start()
		case "stop":
			err = svc.Stop()
		case "backup":
			backup()
		default:
			log.Fatalln("Unknown argument:", flag.Arg(0))
		}
	default:
		log.Fatalln("Unknown arguments:", strings.Join(flag.Args(), " "))
	}
	if err != nil {
		log.Fatalf("failed to %s Metadata: %v", flag.Arg(0), err)
	}
}
