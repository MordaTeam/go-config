package config

import (
	"io"
	"testing/iotest"
)

var _ ConfigProvider = &noProvider{}

// NoProvider is a no-op config provider. It always returns empty io.Reader.
var NoProvider = &noProvider{}

type noProvider struct{}

func (*noProvider) ProvideConfig() (io.Reader, error) {
	return iotest.ErrReader(io.EOF), nil
}
