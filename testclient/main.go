package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dougfort/foods/auth"
	"github.com/dougfort/foods/clienttokens"
)

type config struct {
	makeTokens bool
	tokenSize  int
	tokenCount int
	configPath string
	port       string
}

func main() {
	var cfg config
	var err error

	// parse flags, load config,
	if cfg, err = loadConfig(); err != nil {
		log.Fatalf("loadConfig failed: %s", err)
	}

	if cfg.makeTokens {
		makeTokens(cfg)
	} else {
		runTest(cfg)
	}
}

func makeTokens(cfg config) {
	var outputFile *os.File
	var err error

	log.Printf("creating %d tokens in %s", cfg.tokenCount, cfg.configPath)

	tokens := make([]clienttokens.ClientToken, cfg.tokenCount)
	for i := 0; i < cfg.tokenCount; i++ {
		tokens[i].Client = fmt.Sprintf("client%03d", i+1)
		tokens[i].Token = make([]byte, cfg.tokenSize)
		if _, err = rand.Read(tokens[i].Token); err != nil {
			log.Fatalf("rand.Read failed: %s", err)
		}
	}

	if outputFile, err = os.Create(cfg.configPath); err != nil {
		log.Fatalf("os.Create(%s) failed: %s", cfg.configPath, err)
	}
	defer func() {
		if err = outputFile.Close(); err != nil {
			log.Printf("Close failed: %s", err)
		}
	}()

	encoder := json.NewEncoder(outputFile)
	if err = encoder.Encode(tokens); err != nil {
		log.Printf("Encode failed: %s", err)
	}
}

func runTest(cfg config) {
	var tokens []clienttokens.ClientToken
	var err error

	if tokens, err = clienttokens.Load(cfg.configPath); err != nil {
		log.Fatalf("clienttokens.Load failed: %s", err)
	}

	client := http.Client{}

	t := tokens[0]
	url := fmt.Sprintf("http://localhost:%s/foods/%s?auth=%s",
		cfg.port,
		t.Client,
		auth.String(t.Token, "GET", t.Client, ""))
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("get failed: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("invalid status: (%d) %s",
			resp.StatusCode, resp.Status)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading body: %s", err)
	}

	log.Printf("response: %v", resp)
}

func loadConfig() (config, error) {
	var cfg config

	flag.BoolVar(&cfg.makeTokens, "make-tokens", false, "make tokens for test data")
	flag.IntVar(&cfg.tokenSize, "token-size", 32, "number of bytes in a token")
	flag.IntVar(&cfg.tokenCount, "token-count", 10, "how many tokens to make")
	flag.StringVar(&cfg.configPath, "token-path", "", "path to token config file")
	flag.StringVar(&cfg.port, "port", "8080", "port to contact the server on")

	flag.Parse()

	return cfg, nil
}
