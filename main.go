package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	VERSION = `0.3.1`
)

const (
	BG_COLOR = `#fff`
	FG_COLOR = `#000`
	ENDPOINT = `/hit`
)

var build = `UNKNOWN` // injected in Makefile

var (
	flagBindHostname string
	flagBindPort     int
	flagLogLevel     string
	flagBgColor      string
	flagFgColor      string
	flagEndpoint     string
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
	flag.StringVar(&flagBindHostname, "b", "127.0.0.1", "hostname/ip to bind to")
	flag.IntVar(&flagBindPort, "p", 9999, "port to bind to")
	flag.StringVar(&flagLogLevel, "loglevel", "error", "log level")
	flag.StringVar(&flagEndpoint, "endpoint", ENDPOINT, "endpoint to mount at")
	flag.StringVar(&flagBgColor, "bg", BG_COLOR, "background color, HTML hex string")
	flag.StringVar(&flagFgColor, "fg", FG_COLOR, "foreground color, HTML hex string")
}

func main() {
	flag.Parse()

	setupLogger()
	setupDefaultColors()
	setupDigitsMask()

	go runServerHTTP()
	sigwait()
}
