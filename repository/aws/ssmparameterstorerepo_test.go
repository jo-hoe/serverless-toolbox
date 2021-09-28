package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/jo-hoe/gocommon/repository"
)

var testValue = "testValue"
var testKey = "testKey"
var testPath = "testPath/"

// for mocking suggestion, refer to https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/ssmiface/
// Define a mock struct to be used in your unit tests.
type mockSSM struct {
	ssmiface.SSMAPI
	mapItem map[string]interface{}
	path    string
}

func NewMockSSM(path string, mapItem map[string]interface{}) *mockSSM {
	return &mockSSM{
		mapItem: mapItem,
		path:    path,
	}
}

func createMock() *mockSSM {
	return NewMockSSM(testPath, map[string]interface{}{
		testPath + testKey:       testValue,
		testPath + testKey + "2": testValue + "2",
	})
}

func (mock *mockSSM) PutParameter(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	mock.mapItem[*input.Name] = *input.Value
	return nil, nil
}

func (mock *mockSSM) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	result := new(ssm.GetParameterOutput)
	result.Parameter = new(ssm.Parameter)
	err := errors.New("error")

	if val, ok := mock.mapItem[*input.Name]; ok {
		value := fmt.Sprintf("%v", val)
		result.Parameter.Value = &value
		err = nil
	} else {
		result = nil
	}

	return result, err
}

func (mock *mockSSM) DeleteParameter(input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
	result := new(ssm.DeleteParameterOutput)
	err := errors.New("error")

	if _, ok := mock.mapItem[*input.Name]; ok {
		err = nil
	} else {
		result = nil
	}

	return result, err
}

func (mock *mockSSM) GetParametersByPathPages(input *ssm.GetParametersByPathInput, fn func(*ssm.GetParametersByPathOutput, bool) bool) error {
	err := errors.New("error")

	if *input.Path == mock.path {
		i := 0
		for key, element := range mock.mapItem {
			result := new(ssm.GetParameterOutput)
			result.Parameter = new(ssm.Parameter)
			result.Parameter.Name = &key
			value := fmt.Sprintf("%v", element)
			result.Parameter.Value = &value
			done := fn(&ssm.GetParametersByPathOutput{
				Parameters: []*ssm.Parameter{result.Parameter},
			}, len(mock.mapItem) >= i)
			i++
			if done {
				return nil
			}
		}
		err = nil
	}
	return err
}

func Test_mhm(t *testing.T) {
	session := NewSSMSession("")
	if session == nil {
		t.Error("Session should not be nil")
	}
}

func Test_Find_All(t *testing.T) {
	mock := createMock()
	repo := NewStringSSMParameterStoreRepo(testPath, mock)

	items, err := repo.FindAll()

	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
	if items != nil {
		t.Error("items should not be nil")
	}
	for key, element := range mock.mapItem {
		item := repository.KeyValuePair{
			Key:   key,
			Value: element,
		}
		if !contains(items, item) {
			t.Errorf("Did not find %+v in items slice %+v", item, items)
		}
	}
}

func Test_Find__All_Wrong_Path(t *testing.T) {
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

func Test_Save_And_Find_String(t *testing.T) {
	repo := NewStringSSMParameterStoreRepo(testPath, createMock())
	addedTestKey := "addedTestKey"
	addedTestValue := "addedTestValue"

	repo.Save(addedTestKey, addedTestValue)
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
	addedTestValue := MockItem{
		MockString: "Test",
	}
	mock := mockSSM{
		mapItem: map[string]interface{}{
			testPath + testKey: addedTestValue,
		},
		path: testPath,
	}

	repo := NewSSMParameterStoreRepo(testPath, &mock, MockItem{})

	repo.Save(addedTestKey, addedTestValue)
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
