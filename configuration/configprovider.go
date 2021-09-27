package configuration

import (
	"log"
	"os"
)

// ConfigProvider retrieve configurations content
type ConfigProvider interface {
	GetConfig(configKey string) (string, bool)
	VerboseGetConfig(configKey string) (string, bool)
}

// EnvironmentConfigProvider retrieves configurations from enviroment variables
type EnvironmentConfigProvider struct {
}

// GetConfig returns a configuration for a given key.
// If the variable is present in the environment the
// value (which may be empty) is returned and the boolean is true.
// Otherwise the returned value will be empty and the boolean will
// be false.
func (environmentConfigProvider *EnvironmentConfigProvider) GetConfig(configKey string) (string, bool) {
	return os.LookupEnv(configKey)
}

// VerboseGetConfig behaves like GetConfig but logs error to console if value was not found.
func (environmentConfigProvider *EnvironmentConfigProvider) VerboseGetConfig(configKey string) (string, bool) {
	result, success := environmentConfigProvider.GetConfig(configKey)
	if !success {
		log.Printf("Could not find configuration for key '%v'", configKey)
	}
	return result, success
}
