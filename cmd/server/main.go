package main

import (
	"context"
	"github.com/Hermes-Bird/faraway-test-task.git/config"
	"github.com/Hermes-Bird/faraway-test-task.git/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s := server.NewServer(config.DefaultConfig())

	ctx, cancel := context.WithCancel(context.Background())
	quitCh := make(chan struct{})
	go func() {
		err := s.Start(ctx, quitCh)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	nChan := make(chan os.Signal, 1)

	signal.Notify(nChan, syscall.SIGINT, syscall.SIGTERM)

	<-nChan

	log.Printf("gracefull shutdown...\n")

	cancel()

	select {
	case <-time.NewTimer(config.ServerShutdownTimeout).C:
	case <-quitCh:
	}
}
