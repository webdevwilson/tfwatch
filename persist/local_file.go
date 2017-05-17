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
	path      string
	storeLock *sync.Mutex
	nsLocks   map[string]*sync.Mutex
}

// NewLocalFileStore creates a Store object that stores to the local file system using Glob encoding
func NewLocalFileStore(dir string) (Store, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	l := make(map[string]*sync.Mutex)
	lfs := &localFileStore{dir, &sync.Mutex{}, l}

	// scan for namespaces
	dirs, err := filepath.Glob(path.Join(dir, "*"))
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		ns := path.Base(dir)
		log.Printf("[DEBUG] Found namespace '%s'", ns)
		lfs.CreateNamespace(ns)
	}

	return lfs, err
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

	err = lfs.lock(ns)
	if err != nil {
		return
	}
	defer lfs.unlock(ns)

	guids, err = filepath.Glob(path.Join(lfs.path, ns, "*"))

	if err != nil {
		return
	}

	log.Printf("[DEBUG] Scanning %s directory", lfs.path)
	for i, guid := range guids {
		log.Printf("[DEBUG] Scanning guid '%s'", guid)
		guids[i] = path.Base(guid)
	}
	return
}

// Get returns
func (lfs *localFileStore) Get(ns, guid string, value interface{}) error {

	log.Printf("[DEBUG] Getting item from store namespace: %s guid: %s", ns, guid)
	err := lfs.lock(ns)
	if err != nil {
		return err
	}
	defer lfs.unlock(ns)

	fp := path.Join(lfs.path, ns, guid)

	log.Printf("[DEBUG] Opening %s for reading", fp)
	f, err := os.Open(fp)
	defer func() {
		log.Printf("[DEBUG] Closing file %s", fp)
		f.Close()
	}()

	if err != nil {
		log.Printf("[WARN] Error opening file '%s': %s", fp, err)
		if os.IsNotExist(err) {
			err = &NotFoundError{err}
		}
		return err
	}

	log.Printf("[DEBUG] Decoding value for %s %s", ns, guid)
	return gob.NewDecoder(f).Decode(value)
}

func (lfs *localFileStore) Create(ns string, value interface{}) (string, error) {

	err := lfs.lock(ns)
	if err != nil {
		return "", err
	}
	defer lfs.unlock(ns)

	// create a guid
	guidPtr, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	guid := guidPtr.String()

	// open file
	p := path.Join(lfs.path, ns, guid)
	f, err := os.Create(p)
	defer f.Close()

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

	err := lfs.lock(ns)
	if err != nil {
		return err
	}
	defer lfs.unlock(ns)

	p := path.Join(lfs.path, ns, guid)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("guid '%s' does not exist", guid)
	}

	f, err := os.OpenFile(p, os.O_WRONLY, 0666)
	defer f.Close()

	if err != nil {
		return err
	}
	return gob.NewEncoder(f).Encode(value)
}

func (lfs *localFileStore) Delete(ns, guid string) error {

	err := lfs.lock(ns)
	if err != nil {
		return err
	}
	defer lfs.unlock(ns)

	p := path.Join(lfs.path, ns, guid)
	return os.Remove(p)
}

// CreateNamespace ensures a namespace is configured
func (lfs *localFileStore) CreateNamespace(ns string) error {

	log.Printf("[DEBUG] Creating namespace '%s'", ns)

	lfs.lockStore()
	defer lfs.unlockStore()

	nsPath := path.Join(lfs.path, ns)
	lfs.nsLocks[ns] = &sync.Mutex{}
	return os.MkdirAll(nsPath, os.ModePerm)
}

func (lfs *localFileStore) lock(ns string) error {
	log.Printf("[DEBUG] Locking persist namespace '%s'", ns)

	var v *sync.Mutex
	var ok bool
	if v, ok = lfs.nsLocks[ns]; !ok {
		return fmt.Errorf("Namespace '%s' does not exist!", ns)
	}

	v.Lock()
	return nil
}

func (lfs *localFileStore) unlock(ns string) {
	log.Printf("[DEBUG] Unlocking persist namespace '%s'", ns)
	lfs.nsLocks[ns].Unlock()
}

func (lfs *localFileStore) lockStore() {
	log.Printf("[DEBUG] Locking persist local file store")
	lfs.storeLock.Lock()
}

func (lfs *localFileStore) unlockStore() {
	log.Printf("[DEBUG] Unlocking persist local file store")
	lfs.storeLock.Unlock()
}
