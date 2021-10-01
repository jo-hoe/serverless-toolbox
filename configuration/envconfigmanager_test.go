package configuration

import (
	"os"
	"testing"
)

func TestEnvironmentConfigProviderGetConfig(t *testing.T) {
	testKey := "testKey"
	testValue := "testValue"
	os.Setenv(testKey, testValue)

	provider := EnvironmentConfigProvider{}
	actual, _ := provider.GetConfig(testKey)

	if actual != testValue {
		t.Errorf("Expected %v but retrieved %v", testValue, actual)
	}
}

func TestEnvironmentConfigProviderGetInvalidConfig(t *testing.T) {
	provider := EnvironmentConfigProvider{}
	actual, _ := provider.GetConfig("invalid")

	if actual != nil {
		t.Errorf("Expected empty string but retrieved '%v'", actual)
	}
}

func TestEnvironmentConfigProviderGetEmptyConfig(t *testing.T) {
	provider := EnvironmentConfigProvider{}
	actual, _ := provider.GetConfig("")

	if actual != nil {
		t.Errorf("Expected empty string but retrieved '%v'", actual)
	}
}

func TestEnvironmentConfigProviderResultGetConfig(t *testing.T) {
	testKey := "testKey"
	testValue := "testValue"
	os.Setenv(testKey, testValue)

	provider := EnvironmentConfigProvider{}
	result, err := provider.GetConfig(testKey)

	if err != nil {
		t.Errorf("Expected nil but retrieved %+v", err)
	}
	if result != testValue {
		t.Errorf("Expected %s but got %s", testValue, result)
	}
}

func TestEnvironmentConfigProviderResultGetInvalidConfig(t *testing.T) {
	provider := EnvironmentConfigProvider{}
	result, err := provider.GetConfig("invalid")

	if err == nil {
		t.Errorf("Expected no error but retrieved '%+v'", err)
	}
	if result != nil {
		t.Errorf("Expected nil as result but got '%v'", result)
	}
}

func TestEnvironmentConfigProviderResultGetEmptyConfig(t *testing.T) {
	provider := EnvironmentConfigProvider{}
	_, err := provider.GetConfig("")

	if err == nil {
		t.Errorf("Expected nil but retrieved %+v", err)
	}
}
