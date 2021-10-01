package configuration

// ReadOnlyConfigProvider retrieve configurations content
type ReadOnlyConfigProvider interface {
	GetConfig(configKey string) (string, error)
}

// ConfigReaderWriter retrieve configurations content
type ConfigReaderWriter interface {
	ReadOnlyConfigProvider
	SetConfig(configKey string, configValue interface{}) error
}
