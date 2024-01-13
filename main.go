package main

import (
	"divvy-go-app/internals"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// This channel is used notify various components of the app to notify to that application is closing and thus act accordingly.
	var doneCh = make(chan struct{})

	a := internals.CreateNewApp(doneCh)
	a.Start()
	osCloseCh := make(chan os.Signal, 1)
	signal.Notify(osCloseCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-osCloseCh
	close(doneCh)
}
