package filter

import (
	"testing"

	"github.com/jo-hoe/gocommon/repository"
)

func TestNewPersistantDuplicationFilter(t *testing.T) {
	duplicationFilter := NewPersistantDuplicationFilter(repository.NewInMemoryRepo())

	if duplicationFilter == nil {
		t.Errorf("Notifier was not initialzed")
	}
}

func TestFilter(t *testing.T) {
	duplicationFilter := NewPersistantDuplicationFilter(repository.NewInMemoryRepo())

	count := len(duplicationFilter.Filter(make([]interface{}, 1)))

	if count != 1 {
		t.Errorf("Expected 1 but found %d", count)
	}
}

func TestFilterCurrentDuplicate(t *testing.T) {
	duplicationFilter := NewPersistantDuplicationFilter(repository.NewInMemoryRepo())

	slice := make([]interface{}, 2)

	count := len(duplicationFilter.Filter(slice))

	if count != 1 {
		t.Errorf("Expected 1 but found %d", count)
	}
}

func TestFilterPreviousDuplicate(t *testing.T) {
	duplicationFilter := NewPersistantDuplicationFilter(repository.NewInMemoryRepo())

	duplicationFilter.Filter(make([]interface{}, 1))
	count := len(duplicationFilter.Filter(make([]interface{}, 1)))

	if count != 0 {
		t.Errorf("Expected 0 but found %d", count)
	}
}
