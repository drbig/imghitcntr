package main

import (
	"fmt"
	"net/http"
	"os"
)

import (
	"github.com/sirupsen/logrus"
)

var (
	reqCount = 0
)

func handleRequest(w http.ResponseWriter, req *http.Request) {
	reqCount++ // no locking
	logger.WithFields(logrus.Fields{
		"method":  req.Method,
		"client":  req.RemoteAddr,
		"uri":     req.RequestURI,
		"counter": reqCount,
	}).Info("New request")

	err := req.WriteProxy(os.Stdout)
	if err != nil {
		logger.Error(err)
	}
	fmt.Println("")
	w.Header()["Content-type"] = []string{"text/plain"}
	fmt.Fprintln(w, "Got it")
}
