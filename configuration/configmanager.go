package configuration

// ReadOnlyConfigProvider retrieve configurations content
type ReadOnlyConfigProvider interface {
	GetConfig(configKey string) (string, error)
}

// ConfigReaderManager retrieve configurations content
type ConfigReaderManager interface {
	ReadOnlyConfigProvider
	SetConfig(configKey string, configValue interface{}) error
}
