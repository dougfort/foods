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
func New(dbPath string, names []string) (*Storage, error) {
	var str Storage
	var err error

	if str.Db, err = bolt.Open(dbPath, 0600, nil); err != nil {
		return nil, errors.Wrapf(err, "bold.Open(%s...)", dbPath)
	}

	// create buckets if needed
	for _, name := range names {
		err = str.Db.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return errors.Errorf("create bucket: %s", err)
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "CreateBucketIfNotExists(%s)", name)
		}
	}

	return &str, nil
}

// GetFoods returns the foods a named user likes
func (s *Storage) GetFoods(name string) ([]string, error) {
	var foods []string
	var err error

	err = s.Db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(name))

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			foods = append(foods, string(k))
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "iterate bucket")
	}

	return foods, nil
}

//AddFood stores a new food for the names user
func (s *Storage) AddFood(name string, food string) error {
	var err error

	err = s.Db.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(name))

		err = b.Put([]byte(food), nil)
		return err
	})

	if err != nil {
		return errors.Wrapf(err, "bucket %s put %s", name, food)
	}

	return nil
}

// Close shuts down the database
func (s *Storage) Close() error {
	if err := s.Db.Close(); err != nil {
		return errors.Wrap(err, "Db.Close()")
	}

	return nil
}
