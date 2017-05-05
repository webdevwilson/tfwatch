package persist

import (
	"encoding/gob"
	"log"
	"os"
	"path"
	"path/filepath"

	"sync"

	"fmt"

	uuid "github.com/nu7hatch/gouuid"
)

type localFileStore struct {
	path    string
	lock    *sync.Mutex
	nsLocks map[string]*sync.Mutex
}

// NewLocalFileStore creates a Store object that stores to the local file system using Glob encoding
func NewLocalFileStore(dir string) (store Store, err error) {
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}

	l := make(map[string]*sync.Mutex)
	store = &localFileStore{dir, &sync.Mutex{}, l}
	return
}

// Destroy
func (lfs *localFileStore) Destroy() (err error) {
	l := make(map[string]*sync.Mutex)
	err = os.RemoveAll(lfs.path)
	lfs = &localFileStore{
		lfs.path,
		&sync.Mutex{},
		l,
	}
	return
}

// List returns the keys stored in a namespace
func (lfs *localFileStore) List(ns string) (guids []string, err error) {

	lfs.Lock(ns)
	defer lfs.Unlock(ns)

	guids, err = filepath.Glob(path.Join(lfs.path, ns, "*"))

	if err != nil {
		return
	}

	for i, guid := range guids {
		guids[i] = path.Base(guid)
	}
	return
}

// Get returns
func (lfs *localFileStore) Get(ns, guid string, value interface{}) error {

	lfs.Lock(ns)
	defer lfs.Unlock(ns)

	// ensure the directory exists
	nsPath := path.Join(lfs.path, ns)
	err := os.MkdirAll(nsPath, os.ModePerm)
	if err != nil {
		return err
	}

	fp := path.Join(nsPath, guid)

	f, err := os.Open(fp)

	if err != nil {
		if os.IsNotExist(err) {
			err = &NotFoundError{err}
		}
		return err
	}

	return gob.NewDecoder(f).Decode(value)
}

func (lfs *localFileStore) Create(ns string, value interface{}) (string, error) {

	lfs.Lock(ns)
	defer lfs.Unlock(ns)

	// create a guid
	guidPtr, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	guid := guidPtr.String()

	// open file
	p := path.Join(lfs.path, ns, guid)
	f, err := os.Create(p)
	if err != nil {
		return "", err
	}

	// write to file
	err = gob.NewEncoder(f).Encode(value)
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(p); os.IsNotExist(err) {
		return "", err
	}

	return guid, nil
}

func (lfs *localFileStore) Update(ns, guid string, value interface{}) error {

	if !lfs.HasNamespace(ns) {
		return fmt.Errorf("Namespace %s does not exist", ns)
	}

	lfs.Lock(ns)
	defer lfs.Unlock(ns)

	p := path.Join(lfs.path, ns, guid)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("guid '%s' does not exist", guid)
	}

	f, err := os.OpenFile(p, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	return gob.NewEncoder(f).Encode(value)
}

func (lfs *localFileStore) Delete(ns, guid string) error {

	if !lfs.HasNamespace(ns) {
		return fmt.Errorf("Namespace %s does not exist", ns)
	}

	lfs.Lock(ns)
	defer lfs.Unlock(ns)

	p := path.Join(lfs.path, ns, guid)
	return os.Remove(p)
}

// CreateNamespace ensures a namespace is configured
func (lfs *localFileStore) CreateNamespace(ns string) error {

	log.Printf("[DEBUG] Creating namespace '%s'", ns)

	if lfs.HasNamespace(ns) {
		log.Printf("[DEBUG] Namespace '%s' exists", ns)
		return nil
	}

	lfs.storeLock()
	defer lfs.storeUnlock()

	nsPath := path.Join(lfs.path, ns)
	lfs.nsLocks[ns] = &sync.Mutex{}
	return os.MkdirAll(nsPath, os.ModePerm)
}

func (lfs *localFileStore) HasNamespace(ns string) (ok bool) {
	lfs.storeLock()
	defer lfs.storeUnlock()

	_, ok = lfs.nsLocks[ns]
	return
}

func (lfs *localFileStore) RemoveNamespace(ns string) error {

	log.Printf("[DEBUG] Removing namespace '%s'", ns)

	lfs.storeLock()
	defer lfs.storeUnlock()

	nsPath := path.Join(lfs.path, ns)
	err := os.RemoveAll(nsPath)
	if err != nil {
		lfs.Unlock(ns)
		return err
	}
	delete(lfs.nsLocks, ns)
	return nil
}

func (lfs *localFileStore) Lock(ns string) {
	log.Printf("[DEBUG] Locking persist namespace '%s'", ns)
	lfs.nsLocks[ns].Lock()
}

func (lfs *localFileStore) Unlock(ns string) {
	log.Printf("[DEBUG] Unlocking persist namespace '%s'", ns)
	lfs.nsLocks[ns].Unlock()
}

func (lfs *localFileStore) storeLock() {
	log.Printf("[DEBUG] Locking persist local file store")
	lfs.lock.Lock()
}

func (lfs *localFileStore) storeUnlock() {
	log.Printf("[DEBUG] Unlocking persist local file store")
	lfs.lock.Unlock()
}
