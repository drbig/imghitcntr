// +build !windows

package main

import (
	"os"
	"os/signal"
	"syscall"
)

// sigwait processes signals such as a CTRL-C hit.
func sigwait() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-sig
	logger.Infoln("Signal received, stopping")

	return
}
