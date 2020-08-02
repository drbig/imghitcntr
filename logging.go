package main

import (
	"os"
)

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func setupLogger() {
	level, err := logrus.ParseLevel(flagLogLevel)
	if err != nil {
		panic(err)
	}
	logger = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &logrus.TextFormatter{FullTimestamp: true},
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}
}
