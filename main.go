package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	VERSION = `0.4.0`
)

const (
	BG_COLOR = `#fff`
	FG_COLOR = `#000`
	ENDPOINT = `/hit`
	DB_SIZE  = 32
	DATE_FMT = `2006-01-02`
)

var build = `UNKNOWN` // injected in Makefile

var (
	flagBindHostname string
	flagBindPort     int
	flagLogLevel     string
	flagBgColor      string
	flagFgColor      string
	flagEndpoint     string
	flagCSVPath      string
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
	flag.StringVar(&flagCSVPath, "csv", "", "path to save and load the CSV data dump")
}

func main() {
	flag.Parse()

	setupLogger()
	setupDefaultColors()
	setupDigitsMask()

	if err := loadDB(flagCSVPath); err != nil {
		logger.Errorf("Failed to load state: %s", err)
	}

	go runServerHTTP()
	sigwait()

	if err := saveDB(flagCSVPath); err != nil {
		logger.Errorf("Failed to save state: %s", err)
	}
}
