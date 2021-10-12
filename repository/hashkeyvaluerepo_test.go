package repository

import (
	"reflect"
	"testing"
)

func TestNewHashKeyValueRepo(t *testing.T) {
	repo := NewHashKeyValueRepo(&mockRepo{})

	if repo == nil {
		t.Errorf("Repo was not initialzed")
	}
}

func TestHashKeyValueRepoLength(t *testing.T) {
	repo := NewHashKeyValueRepo(&mockRepo{})

	count, _ := repo.Count()
	if count != 1 {
		t.Errorf("Expected count of 1 but found %d items", count)
	}
}

func TestHashKeyValueRepoToKey(t *testing.T) {
	a := ToKey(mockRepo{test: "dummy"})
	b := ToKey(mockRepo{test: "dummy"})

	if a != b {
		t.Errorf("Expected '%s' be equal to '%s'", a, b)
	}
}

func TestHashKeyValueRepoToKeyUnequal(t *testing.T) {
	a := ToKey(mockRepo{test: "dummy"})
	b := ToKey(mockRepo{test: "dummie"})

	if a == b {
		t.Errorf("Expected '%s' be unequal to '%s'", a, b)
	}
}

func TestHashKeyValueRepoContains(t *testing.T) {
	repo := NewHashKeyValueRepo(&mockRepo{})

	item := "a"
	_, err := repo.Save(item)
	checkError(err, t)

	if !repo.Contains(item) {
		t.Errorf("Could not find %v", item)
	}
}

func TestHashKeyValueRepoOverwrite(t *testing.T) {
	repo := NewHashKeyValueRepo(&mockRepo{})

	item := "a"
	_, err := repo.Overwrite(item)
	checkError(err, t)

	if !repo.Contains(item) {
		t.Errorf("Could not find %v", item)
	}
}

func TestHashKeyValueRepoContainsValue(t *testing.T) {
	repo := NewHashKeyValueRepo(&mockRepo{})

	item := "a"
	_, err := repo.Save(item)
	checkError(err, t)

	if !repo.ContainsValue(item) {
		t.Errorf("Could not find %v", item)
	}
}

type mockRepo struct {
	test string
}

// FindAll calls function of wrapped repository
func (repo *mockRepo) FindAll() ([]KeyValuePair, error) {
	return []KeyValuePair{mockKeyValuePair}, nil
}

// Save calls function of wrapped repository
func (repo *mockRepo) Save(key string, in interface{}) (KeyValuePair, error) {
	return mockKeyValuePair, nil
}

// Save calls function of wrapped repository
func (repo *mockRepo) Overwrite(key string, in interface{}) (KeyValuePair, error) {
	return mockKeyValuePair, nil
}

// Delete calls function of wrapped repository
func (repo *mockRepo) Delete(key string) error {
	return nil
}

// Find calls function of wrapped repository
func (repo *mockRepo) Find(key string) (KeyValuePair, error) {
	return mockKeyValuePair, nil
}

// Find calls function of wrapped repository
func (repo *mockRepo) ContainsValue(in interface{}) bool {
	return reflect.DeepEqual(in, mockKeyValuePair)
}

var mockKeyValuePair = KeyValuePair{
	Key:   "1",
	Value: "test",
}
