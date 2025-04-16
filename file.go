package config

import (
	"fmt"
	"io"
	"os"
)

var _ ConfigProvider = &fileProvider{}

type fileProvider struct {
	cfgPath string
}

// ProvideConfig implements ConfigProvider.
func (f *fileProvider) ProvideConfig() (r io.Reader, err error) {
	file, err := os.Open(f.cfgPath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	return file, nil
}

// FromFile creates a new config provider from a config file.
func FromFile(cfgPath string) *fileProvider {
	return &fileProvider{
		cfgPath: cfgPath,
	}
}
