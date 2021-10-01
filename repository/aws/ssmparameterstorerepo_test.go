package aws

import (
	"testing"

	"github.com/jo-hoe/gocommon/repository"
	"github.com/jo-hoe/gocommon/serialization"
)

var testValue = "testValue"
var testKey = "testKey"
var testPath = "testPath/"

func createMock() *mockSSM {
	return NewMockSSM(testPath, map[string]interface{}{
		testPath + testKey:       testValue,
		testPath + testKey + "2": testValue + "2",
	})
}

func Test_NewSSMSession(t *testing.T) {
	session := NewSSMSession("")
	if session == nil {
		t.Error("Session should not be nil")
	}
}

func Test_Find_All_Strings(t *testing.T) {
	mock := createMock()
	repo := NewStringSSMParameterStoreRepo(testPath, mock)

	items, err := repo.FindAll()

	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
	if items == nil {
		t.Error("items should not be nil")
	}
	if len(items) != len(mock.mapItem) {
		t.Errorf("Expected %d items but found %d", len(mock.mapItem), len(items))
	}
	for key, element := range mock.mapItem {
		item := repository.KeyValuePair{
			Key:   key[:len(testPath)-1],
			Value: element,
		}
		if !contains(items, item) {
			t.Errorf("Did not find %+v in items slice %+v", item, items)
		}
	}
}

func Test_Find_All_Wrong_Path(t *testing.T) {
	repo := NewStringSSMParameterStoreRepo("wrongPath", createMock())

	items, err := repo.FindAll()

	if err == nil {
		t.Error("Error should not be nil")
	}
	if items != nil {
		t.Errorf("Items should be nil but found %+v", items)
	}
}

func Test_Save(t *testing.T) {
	repo := NewStringSSMParameterStoreRepo(testPath, createMock())
	addedTestKey := "addedTestKey"
	addedTestValue := "addedTestValue"

	result, err := repo.Save(addedTestKey, addedTestValue)

	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
	empty := repository.KeyValuePair{}
	if result == empty {
		t.Error("result was empty")
	}
	if result.Key != addedTestKey {
		t.Errorf("Expected %s but received %s", testKey, addedTestKey)
	}
	if result.Value != addedTestValue {
		t.Errorf("Expected %s but received %s", testValue, addedTestValue)
	}
}

func Test_Save_Twice(t *testing.T) {
	repo := NewStringSSMParameterStoreRepo(testPath, createMock())
	addedTestKey := "addedTestKey"
	addedTestValue := "addedTestValue"

	_, err := repo.Save(addedTestKey, addedTestValue)
	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
	result, err := repo.Save(addedTestKey, addedTestValue+"mod")

	if err == nil {
		t.Error("Expected error but found non.")
	}
	empty := repository.KeyValuePair{}
	if result != empty {
		t.Error("Result was empty")
	}
}

func Test_Save_And_Find_String(t *testing.T) {
	repo := NewStringSSMParameterStoreRepo(testPath, createMock())
	addedTestKey := "addedTestKey"
	addedTestValue := "addedTestValue"

	_, err := repo.Save(addedTestKey, addedTestValue)
	if err != nil {
		t.Errorf("Error was not nil %+s", err)
	}

	result, err := repo.Find(addedTestKey)

	if err != nil {
		t.Errorf("Error was not nil %+s", err)
	}
	if result.Value != addedTestValue {
		t.Errorf("Expected %s but found %s", addedTestValue, result.Value)
	}
}

func Test_Save_And_Find(t *testing.T) {
	addedTestKey := "addedTestKey"
	addedTestValue := serialization.MockItem{
		MockString: "Test",
	}
	mock := mockSSM{
		mapItem: map[string]interface{}{
			testPath + testKey: addedTestValue,
		},
		path: testPath,
	}

	repo := NewSSMParameterStoreRepo(testPath, &mock, serialization.MockItem{})

	_, err := repo.Save(addedTestKey, addedTestValue)
	if err != nil {
		t.Errorf("Error was not nil %+s", err)
	}

	result, err := repo.Find(addedTestKey)

	if err != nil {
		t.Errorf("Error was not nil %+s", err)
	}
	if result.Value != addedTestValue {
		t.Errorf("Expected %s but found %s", addedTestValue, result.Value)
	}
}

func Test_Delete(t *testing.T) {
	repo := NewStringSSMParameterStoreRepo(testPath, createMock())

	err := repo.Delete(testKey)

	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
}

func Test_Delete_Wrong_path(t *testing.T) {
	repo := NewStringSSMParameterStoreRepo("wrongPath", createMock())

	err := repo.Delete(testKey)

	if err == nil {
		t.Error("Error should not be nil")
	}
}

func Test_Find(t *testing.T) {
	repo := NewStringSSMParameterStoreRepo(testPath, createMock())

	value, err := repo.Find(testKey)

	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
	if value.Key != testKey {
		t.Errorf("Expected %s but received %s", testKey, value.Key)
	}
	if value.Value != testValue {
		t.Errorf("Expected %s but received %s", testValue, value.Value)
	}
	empty := repository.KeyValuePair{}
	if value == empty {
		t.Error("Value is empty")
	}
}

func Test_Find_Wrong_Path(t *testing.T) {
	repo := NewStringSSMParameterStoreRepo("wrongPath", createMock())

	value, err := repo.Find(testKey)

	if err == nil {
		t.Error("Error should not be nil")
	}
	empty := repository.KeyValuePair{}
	if value != empty {
		t.Errorf("Value should be empty but it is %+s", value)
	}
}

func contains(s []repository.KeyValuePair, e repository.KeyValuePair) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
