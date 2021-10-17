package configuration

// ReadOnlyConfigProvider retrieve configurations content
type ReadOnlyConfigProvider interface {
	GetConfig(configKey string) (interface{}, error)
}

// ConfigManager retrieve configurations content
type ConfigManager interface {
	ReadOnlyConfigProvider
	SetConfig(configKey string, configValue interface{}) error
}
