package main

import (
	"flag"
	"log"
	"runtime"
	"strings"

	"github.com/vharitonsky/iniflags"
)

// OS is the running program's operating system
const OS = runtime.GOOS

var unix, host, port, logPath *string

func main() {
	flag.StringVar(&mc.Server, "dbserver", "localhost", "Metadata Database Server Address")
	flag.IntVar(&mc.Port, "dbport", 27017, "Metadata Database Port")
	flag.StringVar(&mc.Database, "database", "", "Metadata Database Name")
	flag.StringVar(&mc.Collection, "collection", "", "Metadata Database Collection Name")
	flag.StringVar(&mc.Username, "username", "", "Metadata Database Username")
	flag.StringVar(&mc.Password, "password", "", "Metadata Database Password")
	unix = flag.String("unix", "", "UNIX-domain Socket")
	host = flag.String("host", "127.0.0.1", "Server Host")
	port = flag.String("port", "12345", "Server Port")
	logPath = flag.String("log", "", "Log Path")
	iniflags.SetConfigFile("config.ini")
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	switch flag.NArg() {
	case 0:
		run()
	case 1:
		switch flag.Arg(0) {
		case "run":
			run()
		case "backup":
			backup()
		default:
			log.Fatalf("Unknown argument: %s", flag.Arg(0))
		}
	default:
		log.Fatalf("Unknown arguments: %s", strings.Join(flag.Args(), " "))
	}
}
