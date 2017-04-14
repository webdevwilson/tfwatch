package persist

// NotFoundError returned when an object is not found on a Get
type NotFoundError struct {
	error
}

// Store is a simple key-value store that uses namespace-based pessimistic locking
type Store interface {
	// List retrieves the guids in a namespace
	List(ns string) ([]string, error)

	// Get retrieves a value from the Store by guid
	Get(ns, guid string, value interface{}) error

	// Save stores a value, and returns the guid, if any error is returned, nothing is saved
	Save(ns string, value interface{}) (string, error)

	// Delete removes a value from the key-value store
	Delete(ns, guid string) error

	// CreateNamespace ensures that a namespace exists
	CreateNamespace(ns string) error

	// HasNamespace returns whether or not a namespace exists
	HasNamespace(ns string) bool

	// Lock sets a lock on the namespace, when error is returned no lock was acquired
	Lock(ns string)

	// Unlock unlocks a namespace
	Unlock(ns string)
}
