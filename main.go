package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils"
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

	flag.StringVar(&config.Server, "dbserver", "localhost", "Metadata Database Server Address")
	flag.BoolVar(&config.SRV, "srv", false, "SRV Server")
	flag.IntVar(&config.Port, "dbport", 27017, "Metadata Database Port")
	flag.StringVar(&config.Database, "database", "", "Metadata Database Name")
	flag.StringVar(&config.Collection, "collection", "", "Metadata Database Collection Name")
	flag.StringVar(&config.Username, "username", "", "Metadata Database Username")
	flag.StringVar(&config.Password, "password", "", "Metadata Database Password")
	flag.StringVar(&server.Unix, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "127.0.0.1", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	flag.StringVar(&svc.Options.UpdateURL, "update", "", "Update URL")
	exclude := flag.String("exclude", "", "Exclude Files")
	logPath = flag.String("log", "", "Log Path")
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.SetAllowUnknownFlags(true)
	iniflags.Parse()

	svc.Options.ExcludeFiles = strings.Split(*exclude, ",")

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
		case "restart":
			err = svc.Restart()
		case "update":
			err = svc.Update()
		case "backup":
			backup()
		default:
			log.Fatalln("Unknown argument:", flag.Arg(0))
		}
	case 2:
		switch flag.Arg(0) {
		case "restore":
			if utils.Confirm("Do you want to restore database?", 3) {
				restore(flag.Arg(1))
			}
		default:
			log.Fatalln("Unknown arguments:", strings.Join(flag.Args(), " "))
		}
	default:
		log.Fatalln("Unknown arguments:", strings.Join(flag.Args(), " "))
	}
	if err != nil {
		log.Fatalf("failed to %s: %v", flag.Arg(0), err)
	}
}
