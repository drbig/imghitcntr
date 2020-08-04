package main

import (
	"expvar"
	"fmt"
	"image/gif"
	"net/http"
	"os"
)

import (
	"github.com/sirupsen/logrus"
)

const (
	BG_COLOR_KEY = `bg`
	FG_COLOR_KEY = `fg`
	COUNTER_KEY  = `key`
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
	var ok bool

	cntReq.Add(1)
	logger.WithFields(logrus.Fields{
		"method":  req.Method,
		"client":  req.RemoteAddr,
		"counter": cntReq.Value(),
	}).Infof("[%d] New request", cntReq.Value())

	if logger.IsLevelEnabled(logrus.DebugLevel) {
		err := req.WriteProxy(os.Stdout)
		if err != nil {
			logger.Error(err)
		}
		fmt.Println("")
	}

	if err := req.ParseForm(); err != nil {
		cntReqErrors.Add(1)
		logger.Errorf("[%d] Failed to parse form (%d): %s", cntReq.Value(), cntReqErrors.Value(), err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"success\": false, \"msg\": \"failed to parse form\"}\n"))
		return
	}

	key := req.FormValue(COUNTER_KEY)
	if key == "" {
		cntReqErrors.Add(1)
		logger.Warnf("[%d] Request without key (%d)", cntReq.Value(), cntReqErrors.Value())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"success\": false, \"msg\": \"no key present\"}\n"))
		return
	}

	bg := colorBG
	if bgcs := req.FormValue(BG_COLOR_KEY); bgcs != "" {
		logger.Debugf("About to parse BG color param: %s", bgcs)
		bg, ok = parseHexColor(bgcs)
		if !ok {
			cntReqErrors.Add(1)
			logger.Warnf("[%d] Failed to parse BG color (%d): %s", cntReq.Value(), cntReqErrors.Value(), bgcs)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"success\": false, \"msg\": \"failed to parse bg color, use HTML hex string\"}\n"))
			return
		}
	}

	fg := colorFG
	if fgcs := req.FormValue(FG_COLOR_KEY); fgcs != "" {
		logger.Debugf("About to parse FG color param: %s", fgcs)
		fg, ok = parseHexColor(fgcs)
		if !ok {
			cntReqErrors.Add(1)
			logger.Warnf("[%d] Failed to parse FG color (%d): %s", cntReq.Value(), cntReqErrors.Value(), fgcs)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"success\": false, \"msg\": \"failed to parse fg color, use HTML hex string\"}\n"))
			return
		}
	}

	hits := getCount(key)
	logger.Debugf("Generating image for %d hits of %s (%v/%v)...", hits, key, bg, fg)
	img := genImage(hits, bg, fg)

	w.Header().Set("Content-type", "image/gif")
	w.Header().Set("Cache-Control", "no-store")
	err := gif.Encode(w, img, nil)
	if err != nil {
		cntReqErrors.Add(1)
		logger.Errorf("[%d] Failed to encode image (%d): %s", cntReq.Value, cntReqErrors.Value(), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"success\": false, \"msg\": \"failed to encode image, sorry\"}\n"))
	}
}
