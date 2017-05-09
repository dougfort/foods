package main

import (
	"context"
	// TODO: use key/value logging
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dougfort/foods/config"
	"github.com/dougfort/foods/server"
	"github.com/dougfort/foods/storage"
)

func main() {
	os.Exit(run())
}

func run() int {
	var cfg config.Config
	var str storage.Storage
	var err error

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// parse flags, load config,
	if cfg, err = config.Load(); err != nil {
		log.Printf("config.Load failed: %s", err)
		return -1
	}

	// open db,
	if std, err = storage.New(); err != nil {
		log.Printf("storage.New() failed: %s", err)
		return -1
	}

	ctx, cancel := context.WithCancel(context.Background())

	// start server (this will run until the program is killed)
	// TODO: orderly shutdown with cancel
	go foodserv.Serve(ctx, cfg, str)

	// block until sigterm
	s := <-sigChan
	log.Printf("received signal: %s", s)

	cancel()

	return 0
}
