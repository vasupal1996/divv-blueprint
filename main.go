package main

import (
	"context"
	"divvy-go-app/internals"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	/*
		Creating a system level context for proper shutdown behavior when closing the app.
		Every component that requires a graceful shutdown should implement this ctx.
		Eg:
			for run == true {
				select {
				case <-com1.Ctx.Done():
					run = false
					cmp1.Close()
					cmp1.Logger.Debug().Msg("gracefully closed the component cmp1")
				}
			}
	*/
	ctx, cancel := context.WithCancel(context.Background())
	app := internals.CreateNewApp(ctx)
	app.Start()

	/*
		Creating a channel that listens for os level signals for close signals.
		Once signal is detected app starts closing all the resources and components.
	*/
	osCloseCh := make(chan os.Signal, 1)
	signal.Notify(osCloseCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-osCloseCh
	cancel()
	// Waiting for app to close gracefully
	app.Close()
}
