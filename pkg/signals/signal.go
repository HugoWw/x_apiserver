package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var onlyOneSignalHandler = make(chan struct{})

// GetStopSignal registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func GetStopSignal() <-chan struct{} {
	close(onlyOneSignalHandler)
	shutdownQuit := make(chan os.Signal, 2)

	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(shutdownQuit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdownQuit
		cancel()
		<-shutdownQuit

		os.Exit(1) // second signal Exit directly
	}()

	return ctx.Done()
}
