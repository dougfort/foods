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
	const food1 = "broccolli"
	const food2 = "mango"

	var tokens []clienttokens.ClientToken
	var foods []string
	var err error

	if tokens, err = clienttokens.Load(cfg.configPath); err != nil {
		log.Fatalf("clienttokens.Load failed: %s", err)
	}

	baseURL := fmt.Sprintf("http://localhost:%s/foods", cfg.port)

	// verify that the user has no foods at start
	if foods, err = get(baseURL, tokens[0]); err != nil {
		log.Fatalf("first get failed: %s", err)
	}
	if len(foods) != 0 {
		log.Fatalf("expected empty foods: %s", foods)
	}

	// add a food
	if err = post(baseURL, tokens[0], food1); err != nil {
		log.Fatalf("first post failed: %s", err)
	}

	// now we should see a food
	if foods, err = get(baseURL, tokens[0]); err != nil {
		log.Fatalf("2nd get failed: %s", err)
	}
	if len(foods) != 1 {
		log.Fatalf("expected 1 foods: %s", foods)
	}
	if foods[0] != food1 {
		log.Fatalf("invalid food '%s' expected '%s'", foods[0], food1)
	}

	// add another food
	if err = post(baseURL, tokens[0], food2); err != nil {
		log.Fatalf("first post failed: %s", err)
	}

	// now we should see two foods
	if foods, err = get(baseURL, tokens[0]); err != nil {
		log.Fatalf("3rd get failed: %s", err)
	}
	if len(foods) != 2 {
		log.Fatalf("expected 1 foods: %s", foods)
	}
	if foods[1] != food2 {
		log.Fatalf("invalid food '%s' expected '%s'", foods[1], food2)
	}
}

func get(
	baseURL string,
	t clienttokens.ClientToken,
) ([]string, error) {
	var foods []string
	var data []byte
	var err error

	client := http.Client{}
	url := fmt.Sprintf("%s/%s?auth=%s",
		baseURL,
		t.Client,
		auth.String(t.Token, "GET", t.Client, ""))
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Get failed: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: (%d) %s",
			resp.StatusCode, resp.Status)
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %s", err)
	}

	if err = json.Unmarshal(data, &foods); err != nil {
		return nil, fmt.Errorf("Unmarshal failed: %s", err)
	}

	return foods, nil
}

func post(
	baseURL string,
	t clienttokens.ClientToken,
	food string,
) error {
	var err error

	client := http.Client{}
	url := fmt.Sprintf("%s/%s/%s?auth=%s",
		baseURL,
		t.Client,
		food,
		auth.String(t.Token, "POST", t.Client, food))
	resp, err := client.Post(url, "text/html", nil)
	if err != nil {
		return fmt.Errorf("Post failed: %s", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status: (%d) %s",
			resp.StatusCode, resp.Status)
	}

	return nil
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
