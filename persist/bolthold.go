package persist

import (
	"bytes"
	"encoding/gob"
	"log"
	"path"
	"strconv"
	"time"

	"fmt"

	"github.com/boltdb/bolt"
)

type boltStore struct {
	db *bolt.DB
}

// NewBoltStore creates a new store using bolthold
func NewBoltStore(dir string) (Store, error) {
	dbfile := path.Join(dir, "bolt.db")

	log.Printf("[INFO] Bolt DB %s", dbfile)
	db, err := bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 5 * time.Second})

	if err != nil {
		return nil, err
	}

	return &boltStore{db}, nil
}

// List retrieves the guids in a namespace
func (b *boltStore) List(ns string) (guids []string, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(ns))

		return bkt.ForEach(func(k, v []byte) error {
			guids = append(guids, string(k))
			return nil
		})
	})
	return
}

// Get retrieves a value from the Store by guid
func (b *boltStore) Get(ns, guid string, value interface{}) error {
	return b.db.View(func(tx *bolt.Tx) error {
		bytes := tx.Bucket([]byte(ns)).Get([]byte(guid))
		return decode(bytes, value)
	})
}

// Create stores a value, and returns the guid, if any error is returned, nothing is saved
func (b *boltStore) Create(ns string, value interface{}) (idStr string, err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ns))

		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}

		encoded, err := encode(value)
		if err != nil {
			return err
		}

		idStr = strconv.FormatUint(id, 10)
		err = tx.Bucket([]byte(ns)).Put([]byte(idStr), encoded)
		return err
	})
	return
}

// Update updates a stored value, if value does not exist, an error is returned
func (b *boltStore) Update(ns, guid string, value interface{}) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		encoded, err := encode(value)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte(ns)).Put([]byte(guid), encoded)
	})
}

// Delete removes a value from the key-value store
func (b *boltStore) Delete(ns, guid string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(ns)).Delete([]byte(guid))
	})
}

// CreateNamespace ensures that a namespace exists
func (b *boltStore) CreateNamespace(ns string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ns))
		return err
	})
}

// Destroys the store, removing all persisted data
func (b *boltStore) Destroy() error {
	return fmt.Errorf("Destroy not implemented")
}

func encode(value interface{}) ([]byte, error) {
	var buff bytes.Buffer

	en := gob.NewEncoder(&buff)

	err := en.Encode(value)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func decode(data []byte, value interface{}) error {
	var buff bytes.Buffer
	de := gob.NewDecoder(&buff)

	_, err := buff.Write(data)
	if err != nil {
		return err
	}

	err = de.Decode(value)
	if err != nil {
		return err
	}

	return nil
}
