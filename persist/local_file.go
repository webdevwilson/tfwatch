package persist

import (
	"encoding/gob"
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
func (localFileStore *localFileStore) List(ns string) (guids []string, err error) {

	if !localFileStore.HasNamespace(ns) {
		err = fmt.Errorf("Namespace %s does not exist", ns)
		return
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
func (localFileStore *localFileStore) Get(ns, guid string, value interface{}) error {

	if !localFileStore.HasNamespace(ns) {
		return fmt.Errorf("Namespace %s does not exist", ns)
	}

	localFileStore.Lock(ns)
	defer localFileStore.Unlock(ns)

	// ensure the directory exists
	nsPath := path.Join(localFileStore.path, ns)
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

func (localFileStore *localFileStore) Create(ns string, value interface{}) (string, error) {

	if !localFileStore.HasNamespace(ns) {
		return "", fmt.Errorf("Namespace %s does not exist", ns)
	}

	localFileStore.Lock(ns)
	defer localFileStore.Unlock(ns)

	// create a guid
	guidPtr, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	guid := guidPtr.String()

	// open file
	p := path.Join(localFileStore.path, ns, guid)
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

func (localFileStore *localFileStore) Update(ns, guid string, value interface{}) error {

	if !localFileStore.HasNamespace(ns) {
		return fmt.Errorf("Namespace %s does not exist", ns)
	}

	localFileStore.Lock(ns)
	defer localFileStore.Unlock(ns)

	p := path.Join(localFileStore.path, ns, guid)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("guid '%s' does not exist", guid)
	}

	f, err := os.OpenFile(p, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	return gob.NewEncoder(f).Encode(value)
}

func (localFileStore *localFileStore) Delete(ns, guid string) error {

	if !localFileStore.HasNamespace(ns) {
		return fmt.Errorf("Namespace %s does not exist", ns)
	}

	localFileStore.Lock(ns)
	defer localFileStore.Unlock(ns)

	p := path.Join(localFileStore.path, ns, guid)
	return os.Remove(p)
}

// CreateNamespace ensures a namespace is configured
func (localFileStore *localFileStore) CreateNamespace(ns string) error {

	localFileStore.lock.Lock()
	defer localFileStore.lock.Unlock()

	if localFileStore.HasNamespace(ns) {
		return nil
	}

	nsPath := path.Join(localFileStore.path, ns)
	localFileStore.nsLocks[ns] = &sync.Mutex{}
	return os.MkdirAll(nsPath, os.ModePerm)
}

func (localFileStore *localFileStore) HasNamespace(ns string) (ok bool) {
	_, ok = localFileStore.nsLocks[ns]
	return
}

func (localFileStore *localFileStore) RemoveNamespace(ns string) error {
	localFileStore.Lock(ns)
	nsPath := path.Join(localFileStore.path, ns)
	err := os.RemoveAll(nsPath)
	if err != nil {
		localFileStore.Unlock(ns)
		return err
	}
	delete(localFileStore.nsLocks, ns)
	return nil
}

func (localFileStore *localFileStore) Lock(ns string) {
	localFileStore.nsLocks[ns].Lock()
}

func (localFileStore *localFileStore) Unlock(ns string) {
	localFileStore.nsLocks[ns].Unlock()
}
