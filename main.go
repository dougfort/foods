package main

import (
	"context"
	// TODO: use key/value logging
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dougfort/foods/clienttokens"
	"github.com/dougfort/foods/config"
	"github.com/dougfort/foods/server"
	"github.com/dougfort/foods/storage"
)

func main() {
	os.Exit(run())
}

func run() int {
	var cfg config.Config
	var str *storage.Storage
	var tokens []clienttokens.ClientToken
	var names []string
	var err error

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// parse flags, load config,
	if cfg, err = config.Load(); err != nil {
		log.Printf("config.Load failed: %s", err)
		return -1
	}

	if tokens, err = clienttokens.Load(cfg.ConfigPath); err != nil {
		log.Printf("clienttokens.Load failed: %s", err)
	}

	// get the names for initalzing the db
	names = make([]string, len(tokens))
	for i := range tokens {
		names[i] = tokens[i].Client
	}

	// open db, initialize all existing names
	if str, err = storage.New(cfg.DBPath, names); err != nil {
		log.Printf("storage.New() failed: %s", err)
		return -1
	}
	defer func() {
		if err = str.Close(); err != nil {
			log.Printf("error closing storage: %s", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	// start server (this will run until the program is killed)
	// TODO: orderly shutdown with cancel
	go server.Serve(ctx, cfg, str, tokens)

	// block until sigterm
	s := <-sigChan
	log.Printf("received signal: %s", s)

	cancel()

	return 0
}
