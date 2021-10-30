package configuration

import (
	"fmt"

	"github.com/jo-hoe/serverless-toolbox/repository"
)

// RepositoryConfigProvider retrieves configurations from repository.
// Has an interal cache to redurce calls to repositories
type RepositoryConfigProvider struct {
	repo  repository.KeyValueRepo
	cache map[string]interface{}
}

func NewRepositoryConfigProvider(repo repository.KeyValueRepo) *RepositoryConfigProvider {
	return &RepositoryConfigProvider{
		repo:  repo,
		cache: nil,
	}
}

// GetConfig retrieves all configs into cache. Otherwise nil is returned with an error.
func (repositoryConfigProvider *RepositoryConfigProvider) GetConfig(configKey string) (interface{}, error) {
	// init cache
	if repositoryConfigProvider.cache == nil {
		repositoryConfigProvider.cache = make(map[string]interface{})
		// load all items into cache
		keyValuePairs, err := repositoryConfigProvider.repo.FindAll()
		if err != nil {
			return nil, err
		}
		repositoryConfigProvider.cache = make(map[string]interface{})
		for _, keyValuePair := range keyValuePairs {
			repositoryConfigProvider.cache[keyValuePair.Key] = keyValuePair.Value
		}
	}

	if value, ok := repositoryConfigProvider.cache[configKey]; ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("could not find configuration for key '%s' but found %+v", configKey, repositoryConfigProvider.cache)
	}
}

// SetConfig stores a config in string form. Function overwrites existing values
func (repositoryConfigProvider *RepositoryConfigProvider) SetConfig(configKey string, configValue interface{}) error {
	_, err := repositoryConfigProvider.repo.Overwrite(configKey, configValue)
	// if no error was found, item is put into the cache
	if err != nil {
		repositoryConfigProvider.cache[configKey] = configValue
	}
	return err
}
