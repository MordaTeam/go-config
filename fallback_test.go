package config_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/MordaTeam/go-config"
	"github.com/stretchr/testify/require"
)

func TestFallback(t *testing.T) {
	errProvider := &errProvider{}
	expOkCfg := fbConfig{Foo: "bar"}

	t.Run(
		"Normal_OkFirst",
		func(t *testing.T) {
			fb := config.FallbackProvider(okProvider(), errProvider)
			cfg, err := config.New[fbConfig](fb)
			require.NoError(t, err)
			require.Equal(t, expOkCfg, cfg)
		},
	)

	t.Run(
		"Normal_ErrFirst",
		func(t *testing.T) {
			fb := config.FallbackProvider(errProvider, okProvider())
			cfg, err := config.New[fbConfig](fb)
			require.NoError(t, err)
			require.Equal(t, expOkCfg, cfg)
		},
	)

	t.Run(
		"Fail",
		func(t *testing.T) {
			fb := config.FallbackProvider(errProvider, errProvider)
			cfg, err := config.New[fbConfig](fb)
			require.Error(t, err)
			require.Empty(t, cfg)
		},
	)

	t.Run(
		"NilPrInChain",
		func(t *testing.T) {
			fb := config.FallbackProvider(nil, okProvider())
			cfg, err := config.New[fbConfig](fb)
			require.NoError(t, err)
			require.Equal(t, expOkCfg, cfg)
		},
	)
}

type fbConfig struct {
	Foo string `json:"foo"`
}

type errProvider struct{}

func (p *errProvider) ProvideConfig() (io.Reader, error) {
	return nil, errors.New("something went wrong")
}

func okProvider() config.ConfigProvider {
	return config.FromReader(bytes.NewBuffer([]byte(`{"foo": "bar"}`)))
}
