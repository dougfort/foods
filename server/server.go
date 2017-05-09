package server

import (
	"context"
	"encoding/json"
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
	var name string
	var clientToken []byte
	var ok bool

	if reqAuthSlice, ok = request.URL.Query()["auth"]; !ok {
		log.Printf("not auth param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// url.Value is slices of strings
	reqAuth = reqAuthSlice[0]

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
		s.processGET(reqAuth, name, clientToken, w, request)

	case "POST":
		var food string
		if len(splitPath) < 4 {
			log.Printf("invalid path for POST")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		food = splitPath[3]
		s.processPOST(reqAuth, name, clientToken, food, w, request)

	default:
		log.Printf("unhandled method: %s", request.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *server) processGET(
	reqAuth string,
	name string,
	clientToken []byte,
	w http.ResponseWriter,
	request *http.Request,
) {
	var serverAuth string
	var foods []string
	var marshalledFoods []byte
	var err error

	log.Printf("%s: GET", name)
	serverAuth = auth.String(clientToken, "GET", name, "")
	if reqAuth != serverAuth {
		log.Printf("auth mismatch for GET")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if foods, err = s.str.GetFoods(name); err != nil {
		log.Printf("%s: GetFoods() failed: %s", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if marshalledFoods, err = json.Marshal(foods); err != nil {
		log.Printf("%s: json.Marshal failed: %s", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(marshalledFoods); err != nil {
		log.Printf("%s: error writing body: %s", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *server) processPOST(
	reqAuth string,
	name string,
	clientToken []byte,
	food string,
	w http.ResponseWriter,
	request *http.Request,
) {
	var serverAuth string
	var err error

	log.Printf("%s: POST %s", name, food)
	serverAuth = auth.String(clientToken, "POST", name, food)
	if reqAuth != serverAuth {
		log.Printf("auth mismatch for POST")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = s.str.AddFood(name, food); err != nil {
		log.Printf("%s: AddFood failed: %s", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
