package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	//	"github.com/pkg/errors"

	"github.com/dougfort/foods/config"
	"github.com/dougfort/foods/storage"
)

type server struct {
	ctx context.Context
	cfg config.Config
	str *storage.Storage
}

// Serve HTTP
// we could use a more sophisticated multiplexer like Gorilla, but the
// interface is so simple I think this will suffice
func Serve(
	ctx context.Context,
	cfg config.Config,
	str *storage.Storage,
) {
	s := &server{
		ctx: ctx,
		cfg: cfg,
		str: str,
	}

	http.HandleFunc("/foods", s.foodsHandler)

	// ListenAndServe always returns an error
	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
	log.Printf("debug: http.ListenAndServe returned: %s", err)
}

// foodsHandler handles al requests for the /foods route
func (s *server) foodsHandler(w http.ResponseWriter, request *http.Request) {
}
