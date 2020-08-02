package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	VERSION = `0.0.4`
)

var build = `UNKNOWN` // injected in Makefile

var (
	flagBindHostname string
	flagBindPort     int
	flagLogLevel     string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s (options...)
imghitcntr v%s by Piotr S. Staszewski, see LICENSE.txt
binary build by %s

Options:
`, os.Args[0], VERSION, build)
		flag.PrintDefaults()
	}
	flag.StringVar(&flagBindHostname, "l", "127.0.0.1", "hostname/ip to bind to")
	flag.IntVar(&flagBindPort, "p", 9999, "port to bind to")
	flag.StringVar(&flagLogLevel, "loglevel", "error", "log level")
}

func main() {
	flag.Parse()
	setupLogger()
	go runServerHTTP()
	sigwait()
}
