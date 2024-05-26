package cache

import (
	"errors"
	"sync"
)

// Store is a generic object storage and retrieval interface.
type Store interface {
	Add(key string, obj interface{}) error
	Update(key string, obj interface{}) error
	Delete(key string, obj interface{}) error
	List() []interface{}
	Get(key string) (item interface{}, exists bool, err error)
}

// simpleStore is a thread-safe store implementation using sync.Map.
type simpleStore struct {
	items sync.Map
}

// NewSimpleStore creates a new simpleStore.
func NewSimpleStore() Store {
	return &simpleStore{}
}

// Add inserts an object into the store.
func (s *simpleStore) Add(key string, obj interface{}) error {

	if _, loaded := s.items.LoadOrStore(key, obj); loaded {
		return errors.New("object already exists")
	}
	return nil
}

// Update updates an existing object in the store.
func (s *simpleStore) Update(key string, obj interface{}) error {

	if _, loaded := s.items.Load(key); !loaded {
		return errors.New("object does not exist")
	}
	s.items.Store(key, obj)
	return nil
}

// Delete removes an object from the store.
func (s *simpleStore) Delete(key string, obj interface{}) error {

	s.items.Delete(key)
	return nil
}

// List returns a list of all objects in the store.
func (s *simpleStore) List() []interface{} {
	var list []interface{}
	s.items.Range(func(_, value interface{}) bool {
		list = append(list, value)
		return true
	})
	return list
}

// Get retrieves an object from the store.
func (s *simpleStore) Get(key string) (item interface{}, exists bool, err error) {
	item, exists = s.items.Load(key)
	return item, exists, nil
}

// keyFunc generates a unique key for an object.
func (s *simpleStore) keyFunc(obj interface{}) (string, error) {
	// This is a placeholder key function.
	// In a real implementation, you would generate a key based on the object's unique attributes.
	// For example, if obj is a Kubernetes Pod, you might use the namespace and name as the key.
	key, ok := obj.(string)
	if !ok {
		return "", errors.New("failed to generate key for object")
	}
	return key, nil
}
