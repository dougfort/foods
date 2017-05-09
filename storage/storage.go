package storage

import (
	"github.com/pkg/errors"

	"github.com/boltdb/bolt"
)

// Storage holds the database storage capability
type Storage struct {
	Db *bolt.DB
}

// New returns a new storage object
func New(dbPath string) (*Storage, error) {
	var str Storage
	var err error

	if str.Db, err = bolt.Open(dbPath, 0600, nil); err != nil {
		return nil, errors.Wrapf(err, "bold.Open(%s...)", dbPath)
	}

	return &str, nil
}

// Close shuts down the database
func (s *Storage) Close() error {
	if err := s.Db.Close(); err != nil {
		return errors.Wrap(err, "Db.Close()")
	}

	return nil
}
