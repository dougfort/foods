package config

import "flag"

// Config contains external configuration information
type Config struct {
	Port string

	DBPath string

	ConfigPath string
}

// Load the configuration from external sources
func Load() (Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Port, "port", "8080", "port for the server to listen on")
	flag.StringVar(&cfg.DBPath, "dbpath", "/tmp/foods.db", "path to the database")
	flag.StringVar(&cfg.ConfigPath, "token-path", "", "path to client tokens")

	flag.Parse()
	return cfg, nil
}
