package config

import "io"

var _ ConfigProvider = &readerProvider{}

type readerProvider struct {
	r io.Reader
}

// ProvideConfig implements ConfigProvider.
func (r *readerProvider) ProvideConfig() (io.Reader, error) {
	return r.r, nil
}

// FromReader returns a ConfigProvider that reads from r.
func FromReader(r io.Reader) *readerProvider {
	return &readerProvider{r: r}
}
