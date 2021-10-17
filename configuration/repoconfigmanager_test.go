package configuration

import (
	"testing"

	"github.com/jo-hoe/serverless-toolbox/repository"
)

var testKey = "testKey"
var testValue = "testValue"

func TestRepositoryConfigProviderSaveConfig(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	provider := NewRepositoryConfigProvider(repo)
	err := provider.SetConfig(testKey, testValue)

	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
	actual, err := repo.Find(testKey)
	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}

	expected := repository.KeyValuePair{
		Key:   testKey,
		Value: testValue,
	}
	if actual != expected {
		t.Errorf("Expected %v but retrieved %v", expected, actual)
	}
}

func TestRepositoryConfigProviderSaveConfig_Twice(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	provider := NewRepositoryConfigProvider(repo)
	err := provider.SetConfig(testKey, "someValue")
	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
	err = provider.SetConfig(testKey, testValue)
	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}

	actual, err := repo.Find(testKey)
	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}

	expected := repository.KeyValuePair{
		Key:   testKey,
		Value: testValue,
	}
	if actual != expected {
		t.Errorf("Expected %v but retrieved %v", expected, actual)
	}
}

func TestRepositoryConfigProviderGetConfig(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	_, err := repo.Save(testKey, testValue)
	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
	provider := NewRepositoryConfigProvider(repo)
	actual, _ := provider.GetConfig(testKey)

	if actual != testValue {
		t.Errorf("Expected %v but retrieved %v", testValue, actual)
	}
}

func TestRepositoryConfigProviderGetInvalidConfig(t *testing.T) {
	provider := NewRepositoryConfigProvider(repository.NewInMemoryRepo())
	actual, _ := provider.GetConfig("invalid")

	if actual != nil {
		t.Errorf("Expected empty string but retrieved '%v'", actual)
	}
}

func TestRepositoryConfigProviderGetEmptyConfig(t *testing.T) {
	provider := NewRepositoryConfigProvider(repository.NewInMemoryRepo())
	actual, _ := provider.GetConfig("")

	if actual != nil {
		t.Errorf("Expected empty string but retrieved '%v'", actual)
	}
}

func TestRepositoryConfigProviderResultGetConfig(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	_, err := repo.Save(testKey, testValue)
	if err != nil {
		t.Errorf("Expected nil but found error: %+s", err)
	}
	provider := NewRepositoryConfigProvider(repo)
	result, err := provider.GetConfig(testKey)

	if err != nil {
		t.Errorf("Expected nil but retrieved %+v", err)
	}
	if result != testValue {
		t.Errorf("Expected %s but got %s", testValue, result)
	}
}

func TestRepositoryConfigProviderResultGetInvalidConfig(t *testing.T) {
	provider := NewRepositoryConfigProvider(repository.NewInMemoryRepo())
	result, err := provider.GetConfig("invalid")

	if err == nil {
		t.Errorf("Expected no error but retrieved '%+v'", err)
	}
	if result != nil {
		t.Errorf("Expected nil as result but got '%v'", result)
	}
}

func TestRepositoryConfigProviderResultGetEmptyConfig(t *testing.T) {
	provider := NewRepositoryConfigProvider(repository.NewInMemoryRepo())
	_, err := provider.GetConfig("")

	if err == nil {
		t.Errorf("Expected nil but retrieved %+v", err)
	}
}
