package configuration

// ReadOnlyConfigProvider retrieve configurations content
type ReadOnlyConfigProvider interface {
	GetConfig(configKey string) (string, error)
}

// ConfigManager retrieve configurations content
type ConfigManager interface {
	ReadOnlyConfigProvider
	SetConfig(configKey string, configValue interface{}) error
}
