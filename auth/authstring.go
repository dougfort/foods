package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"strings"
)

// String generates an authentication string
// Note that this version is vulnerable to a replay attack
func String(token []byte, method string, name string, food string) string {
	var hasher hash.Hash

	hasher = sha256.New()

	hasher.Write(token)
	hasher.Write([]byte(strings.ToUpper(method)))
	hasher.Write([]byte(name))
	hasher.Write([]byte(food))

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
