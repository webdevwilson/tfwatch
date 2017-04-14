package persist

import (
	"encoding/gob"
	"os"
	"path"
	"path/filepath"

	"sync"

	uuid "github.com/nu7hatch/gouuid"
)

type localFileStore struct {
	path    string
	lock    *sync.Mutex
	nsLocks map[string]*sync.Mutex
}

// NewLocalFileStore creates a Store object that stores to the local file system using Glob encoding
func NewLocalFileStore(p string) (store Store, err error) {
	p = path.Join(p, "data")

	err = os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return
	}

	l := make(map[string]*sync.Mutex)
	store = &localFileStore{p, &sync.Mutex{}, l}
	return
}

// List returns the keys stored in a namespace
func (localFileStore *localFileStore) List(ns string) (guids []string, err error) {

	if !localFileStore.HasNamespace(ns) {
		localFileStore.CreateNamespace(ns)
	}

	// Insure the namespace exists
	if !localFileStore.HasNamespace(ns) {
		guids = nil
		err = NotFoundError{}
		return
	}

	localFileStore.Lock(ns)
	defer localFileStore.Unlock(ns)

	guids, err = filepath.Glob(path.Join(localFileStore.path, ns, "*"))

	if err != nil {
		return
	}

	for i, guid := range guids {
		guids[i] = path.Base(guid)
	}
	return
}

// Get returns
func (localFileStore *localFileStore) Get(ns, guid string, value interface{}) (err error) {

	if !localFileStore.HasNamespace(ns) {
		localFileStore.CreateNamespace(ns)
	}

	localFileStore.Lock(ns)
	defer localFileStore.Unlock(ns)

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
			err = &NotFoundError{err}
		}
		return
	}

	err = gob.NewDecoder(f).Decode(value)

	return
}

func (localFileStore *localFileStore) Save(ns string, value interface{}) (guid string, err error) {

	localFileStore.Lock(ns)
	defer localFileStore.Unlock(ns)

	if !localFileStore.HasNamespace(ns) {
		localFileStore.CreateNamespace(ns)
	}

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

	localFileStore.Lock(ns)
	defer localFileStore.Unlock(ns)

	if !localFileStore.HasNamespace(ns) {
		err = localFileStore.CreateNamespace(ns)
		if err != nil {
			return err
		}
	}

	p := path.Join(localFileStore.path, ns, guid)
	err = os.Remove(p)
	return
}

// CreateNamespace ensures a namespace is configured
func (localFileStore *localFileStore) CreateNamespace(ns string) (err error) {

	localFileStore.lock.Lock()
	defer localFileStore.lock.Unlock()

	if localFileStore.HasNamespace(ns) {
		return
	}

	nsPath := path.Join(localFileStore.path, ns)
	localFileStore.nsLocks[ns] = &sync.Mutex{}
	err = os.MkdirAll(nsPath, os.ModePerm)
	return
}

func (localFileStore *localFileStore) HasNamespace(ns string) (ok bool) {
	_, ok = localFileStore.nsLocks[ns]
	return
}

func (localFileStore *localFileStore) Lock(ns string) {
	localFileStore.nsLocks[ns].Lock()
}

func (localFileStore *localFileStore) Unlock(ns string) {
	localFileStore.nsLocks[ns].Unlock()
}
