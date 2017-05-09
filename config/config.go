package config

import "flag"

// Config contains external configuration information
type Config struct {
	Port string

	// Tokens is the list of valid client tokens
	Tokens []string
}

// Load the configuration from external sources
func Load() (Config, error) {
	var cfg Config

	flag.String(cfg.Port, "8080", "port for the server to listen on")

	return cfg, nil
}
