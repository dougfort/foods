package clienttokens

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"
)

// ClientToken associates a token with an individual client
type ClientToken struct {
	Client string
	Token  []byte
}

// Load loads client tokens from a JSON config file
func Load(configPath string) ([]ClientToken, error) {
	var inputFile *os.File
	var tokens []ClientToken
	var err error

	if inputFile, err = os.Open(configPath); err != nil {
		return nil, errors.Wrapf(err, "os.Open(%s) failed", configPath)
	}
	defer func() {
		if err = inputFile.Close(); err != nil {
			log.Printf("Close failed: %s", err)
		}
	}()

	decoder := json.NewDecoder(inputFile)
	if err = decoder.Decode(&tokens); err != nil {
		return nil, errors.Wrap(err, "Decode failed")
	}

	return tokens, nil
}
