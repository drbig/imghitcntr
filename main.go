package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

const (
	VERSION = "0.0.1"
)

var (
	flagBindHostname string
	flagBindPort     int
	flagLogLevel     string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s (options...)
imghitcntr v%s by Piotr S. Staszewski, see LICENSE.txt

Options:
`, os.Args[0], VERSION)
		flag.PrintDefaults()
	}
	flag.StringVar(&flagBindHostname, "l", "127.0.0.1", "hostname/ip to bind to")
	flag.IntVar(&flagBindPort, "p", 9999, "port to bind to")
	flag.StringVar(&flagLogLevel, "loglevel", "error", "log level")
}

func main() {
	flag.Parse()
	setupLogger()

	bind_addr := fmt.Sprintf("%s:%d", flagBindHostname, flagBindPort)

	http.HandleFunc("/hit", handleRequest)
	logger.Infof("Bind to %s", bind_addr)
	http.ListenAndServe(bind_addr, nil)
}
