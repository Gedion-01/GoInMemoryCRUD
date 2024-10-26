package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Gedion-01/Go-Crud-Challenge/server"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	srv := server.NewServer()

	select {
	case <-stop:
		srv.Stop()
	}
}