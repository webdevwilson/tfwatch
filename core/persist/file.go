package persist

import (
	"encoding/gob"
	"os"
	"path"
	"path/filepath"

	uuid "github.com/nu7hatch/gouuid"
)

// NotFoundError returned when an object is not found on a Get
type NotFoundError struct {
	error
}

// Store is used to store data
type Store interface {
	// List retrieves the guids in a namespace
	List(ns string) ([]string, error)

	// Get retrieves a value from the Store by guid
	Get(ns, guid string, value interface{}) error

	// Save stores a value, and returns the guid, if any error is returned, nothing is saved
	Save(ns string, value interface{}) (string, error)

	Delete(ns, guid string) error

	// CreateNamespace ensures that a namespace exists
	CreateNamespace(ns string) error
}

type localFileStore struct {
	path string
}

// NewLocalFileStore creates a Store object that stores to the local file system using Glob encoding
func NewLocalFileStore(p string) (store Store, err error) {
	p = path.Join(p, "data")

	err = os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return
	}

	store = &localFileStore{p}
	return
}

func (localFileStore *localFileStore) List(ns string) (guids []string, err error) {
	guids, err = filepath.Glob(path.Join(localFileStore.path, ns, "*"))

	if err != nil {
		return
	}

	for i, guid := range guids {
		guids[i] = path.Base(guid)
	}
	return
}

func (localFileStore *localFileStore) Get(ns, guid string, value interface{}) (err error) {

	// ensure the directory exists
	nsPath := path.Join(localFileStore.path, ns)
	err = os.MkdirAll(nsPath, os.ModePerm)
	if err != nil {
		return
	}

	fp := path.Join(nsPath, guid)

	f, err := os.Open(fp)

	if err != nil {
		if os.IsNotExist(err) {
			err = &NotFoundError{}
		}
		return
	}

	err = gob.NewDecoder(f).Decode(value)

	return
}

func (localFileStore *localFileStore) Save(ns string, value interface{}) (guid string, err error) {

	// create a guid
	var guidPtr *uuid.UUID
	guidPtr, err = uuid.NewV4()
	if err != nil {
		return
	}
	guid = guidPtr.String()

	// open file
	p := path.Join(localFileStore.path, ns, guid)
	f, err := os.Create(p)
	if err != nil {
		return
	}

	// write to file
	err = gob.NewEncoder(f).Encode(value)
	if err != nil {
		return
	}

	return
}

func (localFileStore *localFileStore) Delete(ns, guid string) (err error) {
	p := path.Join(localFileStore.path, ns, guid)
	err = os.Remove(p)
	return
}

// CreateNamespace ensures a namespace is configured
func (localFileStore *localFileStore) CreateNamespace(ns string) (err error) {
	nsPath := path.Join(localFileStore.path, ns)
	err = os.MkdirAll(nsPath, os.ModePerm)
	return
}
