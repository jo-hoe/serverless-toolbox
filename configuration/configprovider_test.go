package configuration

import (
	"bytes"
	"log"
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

	if actual != "" {
		t.Errorf("Expected empty string but retrieved '%v'", actual)
	}
}

func TestEnvironmentConfigProviderGetEmptyConfig(t *testing.T) {
	provider := EnvironmentConfigProvider{}
	actual, _ := provider.GetConfig("")

	if actual != "" {
		t.Errorf("Expected empty string but retrieved '%v'", actual)
	}
}

func TestEnvironmentConfigProviderResultGetConfig(t *testing.T) {
	testKey := "testKey"
	testValue := "testValue"
	os.Setenv(testKey, testValue)

	provider := EnvironmentConfigProvider{}
	_, result := provider.GetConfig(testKey)

	if result != true {
		t.Errorf("Expected true but retrieved %v", result)
	}
}

func TestEnvironmentConfigProviderResultGetInvalidConfig(t *testing.T) {
	provider := EnvironmentConfigProvider{}
	_, result := provider.GetConfig("invalid")

	if result != false {
		t.Errorf("Expected true but retrieved %v", result)
	}
}

func TestEnvironmentConfigProviderResultGetEmptyConfig(t *testing.T) {
	provider := EnvironmentConfigProvider{}
	_, result := provider.GetConfig("")

	if result != false {
		t.Errorf("Expected true but retrieved %v", result)
	}
}

func TestEnvironmentVerboseGetConfig(t *testing.T) {
	testKey := "testKey"
	testValue := "testValue"
	os.Setenv(testKey, testValue)

	provider := EnvironmentConfigProvider{}
	actual, _ := provider.VerboseGetConfig(testKey)

	if actual != testValue {
		t.Errorf("Expected %v but retrieved %v", testValue, actual)
	}
}

func TestEnvironmentVerboseGetConfigLog(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	provider := EnvironmentConfigProvider{}
	provider.VerboseGetConfig("")

	logOutput := buf.String()
	if len(logOutput) == 0 {
		t.Error("No log output found")
	}
}
