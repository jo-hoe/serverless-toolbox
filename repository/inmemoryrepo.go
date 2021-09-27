package repository

import (
	"fmt"
	"sync"
)

// InMemoryRepo stores all entities in-memory
type InMemoryRepo struct {
	mapStore map[string]interface{}
	mutex    sync.RWMutex
}

// NewInMemoryRepo creates a new instance of the repository
func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		mapStore: make(map[string]interface{}),
	}
}

// Save one item
func (repo *InMemoryRepo) Save(key string, in interface{}) (KeyValuePair, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	result := KeyValuePair{}
	_, ok := repo.mapStore[key]
	if ok {
		return result, fmt.Errorf("key %s already exists", key)
	}
	repo.mapStore[key] = in
	result.Key = key
	result.Value = in

	return result, nil
}

// FindAll items
func (repo *InMemoryRepo) FindAll() ([]KeyValuePair, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	result := make([]KeyValuePair, 0, len(repo.mapStore))

	for key, value := range repo.mapStore {
		result = append(result, KeyValuePair{
			Key:   key,
			Value: value,
		})
	}

	return result, nil
}

// Delete an item from the repository
func (repo *InMemoryRepo) Delete(key string) error {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	_, ok := repo.mapStore[key]
	if !ok {
		return fmt.Errorf("key %s not found", key)
	}
	delete(repo.mapStore, key)
	return nil
}

// Find retrieves and item from the repository
func (repo *InMemoryRepo) Find(key string) (KeyValuePair, error) {
	result := KeyValuePair{
		Key:   key,
		Value: nil,
	}
	value, ok := repo.mapStore[key]
	if !ok {
		return result, fmt.Errorf("key %s not found", key)
	}
	result.Value = value
	return result, nil
}
