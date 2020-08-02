package main

import (
	"expvar"
	"fmt"
	"net/http"
	"os"
)

import (
	"github.com/sirupsen/logrus"
)

var (
	cntReq       = expvar.NewInt("statsRequests")
	cntReqErrors = expvar.NewInt("statsReqErrors")
)

func runServerHTTP() {
	bind_addr := fmt.Sprintf("%s:%d", flagBindHostname, flagBindPort)
	http.HandleFunc("/hit", handleRequest)
	logger.Infof("Bind to %s", bind_addr)
	http.ListenAndServe(bind_addr, nil)
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	cntReq.Add(1)
	referrer := req.Header.Get("Referer")
	logger.WithFields(logrus.Fields{
		"method":   req.Method,
		"referrer": referrer,
		"client":   req.RemoteAddr,
		"counter":  cntReq.Value(),
	}).Infof("[%d] New request", cntReq.Value())

	if logger.IsLevelEnabled(logrus.DebugLevel) {
		err := req.WriteProxy(os.Stdout)
		if err != nil {
			logger.Error(err)
		}
		fmt.Println("")
	}

	if referrer == "" {
		cntReqErrors.Add(1)
		logger.Warnf("[%d] Request without referrer (%d)", cntReq.Value(), cntReqErrors.Value())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "msg": "no referer header present"}`))
		return
	}

	w.Header()["Content-type"] = []string{"text/plain"}
	fmt.Fprintf(w, "%d", getCount(referrer))
}
