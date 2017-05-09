package storage

import "github.com/boltdb/bolt"

// Storage holds the database storage capability
type Storage struct {
	Db *bolt.DB
}

// New returns a new storage object
func New() (*Storage, error) {
	var str Storage

	return &str, nil
}
