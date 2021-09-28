package aws

import (
	"errors"
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
	mapItem map[string]string
	path    string
}

func NewMockSSM(path string, mapItem map[string]string) *mockSSM {
	return &mockSSM{
		mapItem: mapItem,
		path:    path,
	}
}

func createMock() *mockSSM {
	return NewMockSSM(testPath, map[string]string{
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
		result.Parameter.Value = &val
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
