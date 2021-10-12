package repository

import (
	"testing"
)

type mockStruct struct {
	mockString string
}

var mockInstance = mockStruct{
	mockString: "mock",
}

func TestNewInMemoryRepo(t *testing.T) {
	repo := NewInMemoryRepo()

	if repo == nil {
		t.Errorf("Repo was not initialzed")
	}
}

func TestInMemoryRepoFindAllEmpty(t *testing.T) {
	repo := NewInMemoryRepo()

	items, _ := repo.FindAll()
	if items == nil {
		t.Errorf("Could not get all items of empty repo")
	}
}

func TestInMemoryRepoSave(t *testing.T) {
	repo := NewInMemoryRepo()

	_, err := repo.Save("some key", mockInstance)
	checkError(err, t)

	allItems, _ := repo.FindAll()
	count := len(allItems)
	if count != 1 {
		t.Errorf("Expected 1 item but found %d items", count)
	}
}

func TestInMemoryRepoSaveValue(t *testing.T) {
	repo := NewInMemoryRepo()

	result, _ := repo.Save("some key", mockInstance)

	if result.Value != mockInstance {
		t.Errorf("Expected %v item but found %v", mockInstance, result)
	}
}

func TestInMemoryRepoSaveTwiceError(t *testing.T) {
	repo := NewInMemoryRepo()

	_, err := repo.Save("samekey", mockInstance)
	checkError(err, t)
	_, err = repo.Save("samekey", mockInstance)

	if err == nil {
		t.Error("Expected error was nil although same key was inserted twice")
	}
}

func TestInMemoryRepoOverwrite(t *testing.T) {
	repo := NewInMemoryRepo()

	_, err := repo.Overwrite("samekey", mockInstance)
	checkError(err, t)
	item, err := repo.Overwrite("samekey", mockInstance)

	if err != nil {
		t.Errorf("did not expect an error %v", err)
	}
	if item.Value != mockInstance {
		t.Errorf("expeted %+v but received %+v", mockInstance, err)
	}
}

func TestInMemoryRepoSaveTwiceLength(t *testing.T) {
	repo := NewInMemoryRepo()

	_, err := repo.Save("samekey", mockInstance)
	checkError(err, t)
	_, err = repo.Save("samekey", mockInstance)
	checkFailure(err, t)

	items, _ := repo.FindAll()
	count := len(items)
	if count != 1 {
		t.Errorf("Expected length to be 1 but was %d", count)
	}
}

func TestInMemoryRepoDelete(t *testing.T) {
	repo := NewInMemoryRepo()

	result, _ := repo.Save("some key", mockInstance)
	err := repo.Delete(result.Key)
	checkError(err, t)

	allItems, _ := repo.FindAll()
	count := len(allItems)
	if count != 0 {
		t.Errorf("Expected 0 item but found %d items", count)
	}
}

func TestInMemoryRepoDeleteInvalid(t *testing.T) {
	repo := NewInMemoryRepo()

	err := repo.Delete("invalid")

	if err == nil {
		t.Errorf("Error is nil although Delete was called with an invalid value")
	}
}

func TestInMemoryRepoFind(t *testing.T) {
	repo := NewInMemoryRepo()

	storedItem, _ := repo.Save("some key", mockInstance)
	result, _ := repo.Find(storedItem.Key)

	if result.Value != mockInstance {
		t.Errorf("Expected %v but retrieved %v", mockInstance, result.Value)
	}
}

func TestInMemoryRepoFindInvalid(t *testing.T) {
	repo := NewInMemoryRepo()

	_, err := repo.Find("invalid")

	if err == nil {
		t.Errorf("Error is nil although Find was called with an invalid value")
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func checkFailure(err error, t *testing.T) {
	if err == nil {
		t.Error(err)
	}
}
