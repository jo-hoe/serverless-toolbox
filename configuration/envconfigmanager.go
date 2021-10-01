package configuration

import (
	"fmt"
	"os"
)

// EnvironmentConfigProvider retrieves configurations from enviroment variables
type EnvironmentConfigProvider struct {
}

// GetConfig returns a configuration for a given key. Otherwise nil is returned with an error
func (environmentConfigProvider *EnvironmentConfigProvider) GetConfig(configKey string) (interface{}, error) {
	result, success := os.LookupEnv(configKey)
	if !success {
		return nil, fmt.Errorf("could not find configuration for key '%s'", configKey)
	}
	return result, nil
}
