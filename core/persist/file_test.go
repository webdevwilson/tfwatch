package persist

import (
	"log"
	"os"
	"path"
	"testing"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
)

type data struct {
	Name  string
	Count int
}

var testData = [3]*data{
	&data{"a", 0},
	&data{"b", 1},
	&data{"c", 2},
}

// TestMain sets up a datastore
func createStore() (*localFileStore, func()) {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("[ERROR] Error getting working directory: %s", err)
		os.Exit(1)
	}

	datastorePath := path.Join(wd, ".test")
	store, err := NewLocalFileStore(datastorePath)

	if err != nil {
		log.Printf("[ERROR] Error creating NewLocalFileStore: %s", err)
		os.Exit(1)
	}

	// return the store and cleanup func
	return store.(*localFileStore), func() {
		os.RemoveAll(datastorePath)

		// insure we delete the file
		if _, err := os.Stat(datastorePath); err == nil {
			log.Printf("[ERROR] Error cleaning up unit tests")
		}
	}
}

func Test_CreateNamespace(t *testing.T) {
	store, cleanup := createStore()
	defer cleanup()

	err := store.CreateNamespace(t.Name())
	if err != nil {
		t.Error(err)
	}
}

func Test_Save(t *testing.T) {
	store, cleanup := createStore()
	defer cleanup()

	err := store.CreateNamespace(t.Name())
	if err != nil {
		t.Error(err)
	}

	guid, err := store.Save(t.Name(), &data{"name", 12})
	if err != nil {
		t.Error(err)
	}

	_, err = uuid.ParseHex(guid)

	if err != nil {
		t.Errorf("Invalid guid '%s': %s", guid, err)
	}
}

func Test_List(t *testing.T) {
	store, cleanup := createStore()
	defer cleanup()

	err := store.CreateNamespace(t.Name())
	if err != nil {
		t.Error(err)
	}

	guids := make([]string, len(testData))

	for i, data := range testData {
		guid, err := store.Save(t.Name(), data)
		if err != nil {
			t.Error(err)
		}
		guids[i] = guid
	}

	test, err := store.List(t.Name())

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(testData), len(test))
}

func Test_List_When_Namespace_Doesnt_Exist(t *testing.T) {
	store, cleanup := createStore()
	defer cleanup()

	test, err := store.List(t.Name())

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 0, len(test))
}

func Test_Get(t *testing.T) {
	store, cleanup := createStore()
	defer cleanup()

	err := store.CreateNamespace(t.Name())
	if err != nil {
		t.Error(err)
	}

	// persist test data
	guids := make([]string, len(testData))
	for i, data := range testData {
		guid, err := store.Save(t.Name(), data)
		if err != nil {
			t.Error(err)
		}
		guids[i] = guid
	}

	for i, guid := range guids {
		var data data
		err := store.Get(t.Name(), guid, &data)
		if err != nil {
			t.Error(err)
		}

		expected := testData[i]
		assert.Equal(t, expected.Name, data.Name)
		assert.Equal(t, expected.Count, data.Count)
	}
}

func Test_Get_NonExisting_returns_nil(t *testing.T) {
	store, cleanup := createStore()
	defer cleanup()

	err := store.CreateNamespace(t.Name())
	if err != nil {
		t.Error(err)
	}

	var value data
	err = store.Get(t.Name(), "NON-EXISTING", &value)

	_, ok := err.(*NotFoundError)
	assert.Equal(t, true, ok, "Expected NotFoundError, Got: %s", err)
}

func Test_Get_NonExisting_returns_nil_When_Namespace_Doesnt_Exist(t *testing.T) {
	store, cleanup := createStore()
	defer cleanup()

	var value data
	err := store.Get(t.Name(), "NON-EXISTING", value)

	_, ok := err.(*NotFoundError)
	assert.Equal(t, true, ok)
}

func Test_Delete(t *testing.T) {
	store, cleanup := createStore()
	defer cleanup()

	err := store.CreateNamespace(t.Name())
	if err != nil {
		t.Error(err)
	}

	guid, err := store.Save(t.Name(), &data{"name", 12})
	if err != nil {
		t.Error(err)
	}

	err = store.Delete(t.Name(), guid)

	if err != nil {
		t.Error(err)
	}

	var data data
	err = store.Get(t.Name(), guid, &data)
	if err != nil && !os.IsNotExist(err) {
		t.Error(err)
	}
}
