package configuration

import (
	"github.com/jo-hoe/gocommon/repository"
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
	if err == nil {
		return result.Value, err
	} else {
		return nil, err
	}
}

// SetConfig stores a config in string form.
func (repositoryConfigProvider *RepositoryConfigProvider) SetConfig(configKey string, configValue interface{}) error {
	_, err := repositoryConfigProvider.repo.Save(configKey, configValue)
	return err
}
