package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	//	"github.com/pkg/errors"

	"github.com/dougfort/foods/auth"
	"github.com/dougfort/foods/clienttokens"
	"github.com/dougfort/foods/config"
	"github.com/dougfort/foods/storage"
)

type server struct {
	ctx       context.Context
	cfg       config.Config
	str       *storage.Storage
	clientMap map[string][]byte
}

// Serve HTTP
// we could use a more sophisticated multiplexer like Gorilla, but the
// interface is so simple I think this will suffice
func Serve(
	ctx context.Context,
	cfg config.Config,
	str *storage.Storage,
	tokens []clienttokens.ClientToken,
) {
	s := &server{
		ctx: ctx,
		cfg: cfg,
		str: str,
	}
	s.clientMap = make(map[string][]byte)
	for _, t := range tokens {
		s.clientMap[t.Client] = t.Token
	}

	http.HandleFunc("/foods/", s.foodsHandler)

	// ListenAndServe always returns an error
	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
	log.Printf("debug: http.ListenAndServe returned: %s", err)
}

// foodsHandler handles al requests for the /foods route
func (s *server) foodsHandler(w http.ResponseWriter, request *http.Request) {
	var reqAuthSlice []string
	var reqAuth string
	var splitPath []string
	var method string
	var name string
	var clientToken []byte
	var serverAuth string
	var ok bool

	if reqAuthSlice, ok = request.URL.Query()["auth"]; !ok {
		log.Printf("not auth param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// url.Value is slices of strings
	reqAuth = reqAuthSlice[0]

	method = strings.ToUpper(request.Method)

	// we expect a path like "/foods/<name>/..."
	// this gives us splitPath = ["", "foods", <name>, ...]
	splitPath = strings.Split(request.URL.Path, "/")
	if len(splitPath) < 3 {
		log.Printf("unparseable path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name = splitPath[2]

	if clientToken, ok = s.clientMap[name]; !ok {
		log.Printf("unknown user: '%s'", name)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch strings.ToUpper(request.Method) {
	case "GET":
		serverAuth = auth.String(clientToken, method, name, "")
		if reqAuth != serverAuth {
			log.Printf("auth mismaptch for GET")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	case "POST":
		var food string
		if len(splitPath) < 4 {
			log.Printf("invalid path for POST")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		food = splitPath[3]
		serverAuth = auth.String(clientToken, method, name, food)
		if reqAuth != serverAuth {
			log.Printf("auth mismaptch for POST")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	log.Printf("request path: %s, auth = %s", request.URL.Path, reqAuth)
}
