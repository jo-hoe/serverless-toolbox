package configuration

import (
	"fmt"

	"github.com/jo-hoe/serverless-toolbox/repository"
)

// RepositoryConfigProvider retrieves configurations from enviroment variables
type RepositoryConfigProvider struct {
	repo repository.KeyValueRepo
}

func NewRepositoryConfigProvider(repo repository.KeyValueRepo) *RepositoryConfigProvider {
	return &RepositoryConfigProvider{
		repo: repo,
	}
}

// GetConfig returns a configuration for a given key. Otherwise nil is returned with an error
func (repositoryConfigProvider *RepositoryConfigProvider) GetConfig(configKey string) (interface{}, error) {
	result, err := repositoryConfigProvider.repo.Find(configKey)
	value, erri := repositoryConfigProvider.repo.FindAll()
	fmt.Printf("all: '%v'; error: %v", value, erri)
	if err == nil {
		return result.Value, err
	} else {
		return nil, fmt.Errorf("could not find configuration for key '%s' error: '%v'", configKey, err)
	}
}

// SetConfig stores a config in string form. Function overwrites existing values
func (repositoryConfigProvider *RepositoryConfigProvider) SetConfig(configKey string, configValue interface{}) error {
	_, err := repositoryConfigProvider.repo.Overwrite(configKey, configValue)
	return err
}
