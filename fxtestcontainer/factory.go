package fxtestcontainer

import (
	"fmt"
	"strings"

	"github.com/ankorstore/yokai/config"
)

var _ ContainerConfigFactory = (*DefaultContainerConfigFactory)(nil)

// ContainerConfigFactory is the interface for creating container configurations.
type ContainerConfigFactory interface {
	Create(configKey string) (*ContainerConfig, error)
}

// DefaultContainerConfigFactory is the default ContainerConfigFactory implementation.
type DefaultContainerConfigFactory struct {
	config *config.Config
}

// NewDefaultContainerConfigFactory returns a new DefaultContainerConfigFactory instance.
func NewDefaultContainerConfigFactory(config *config.Config) *DefaultContainerConfigFactory {
	return &DefaultContainerConfigFactory{
		config: config,
	}
}

// Create creates a ContainerConfig from the configuration file.
func (f *DefaultContainerConfigFactory) Create(configKey string) (*ContainerConfig, error) {
	configPath := fmt.Sprintf("modules.testcontainer.containers.%s", configKey)

	if !f.config.IsSet(configPath) {
		return nil, fmt.Errorf("%w: %s", ErrContainerConfigNotFound, configKey)
	}

	// Get environment variables and convert keys back to uppercase
	envVars := f.config.GetStringMapString(fmt.Sprintf("%s.environment", configPath))
	environment := make(map[string]string)
	for key, value := range envVars {
		// Convert lowercase keys back to uppercase for Docker environment variables
		upperKey := strings.ToUpper(key)
		environment[upperKey] = value
	}

	containerConfig := &ContainerConfig{
		Name:         f.config.GetString(fmt.Sprintf("%s.name", configPath)),
		Image:        f.config.GetString(fmt.Sprintf("%s.image", configPath)),
		Port:         f.config.GetString(fmt.Sprintf("%s.port", configPath)),
		ExposedPorts: f.config.GetStringSlice(fmt.Sprintf("%s.exposed_ports", configPath)),
		Environment:  environment,
		Cmd:          f.config.GetStringSlice(fmt.Sprintf("%s.cmd", configPath)),
	}

	// Set default name if not provided
	if containerConfig.Name == "" {
		containerConfig.Name = configKey
	}

	return containerConfig, nil
}
