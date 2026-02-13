package config

import (
	"errors"
	"fmt"
	"io"
)

var _ ConfigProvider = &fallbackProvider{}

type fallbackProvider struct {
	providers []ConfigProvider
}

// Creates new fallback line of providers.
// Call for ProvideConfig() returns first successfull result of ProvideConfig() call of internal provider.
// If all internal providers fail, return resulting error.
func FallbackProvider(prs ...ConfigProvider) ConfigProvider {
	return &fallbackProvider{
		providers: prs,
	}
}

func (p *fallbackProvider) ProvideConfig() (io.Reader, error) {
	var resErr error
	for idx, pr := range p.providers {
		r, err := pr.ProvideConfig()
		if err != nil {
			resErr = errors.Join(resErr, fmt.Errorf("provider %d: %w", idx, err))
			continue
		}

		return r, nil
	}

	return nil, resErr
}
