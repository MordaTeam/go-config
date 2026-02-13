package config_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/MordaTeam/go-config"
	"github.com/stretchr/testify/require"
)

func TestFallbackDecoder(t *testing.T) {
	errDec := &errDecoder{}
	expOkCfg := fbConfig{Foo: "bar"}

	t.Run(
		"Normal_OkFirst",
		func(t *testing.T) {
			fb := config.FallbackDecoder(okDecoder(reader()), errDec)(nil)
			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.NoError(t, err)
			require.Equal(t, expOkCfg, cfg)
		},
	)

	t.Run(
		"Normal_ErrFirst",
		func(t *testing.T) {
			fb := config.FallbackDecoder(errDec, okDecoder(reader()))(nil)
			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.NoError(t, err)
			require.Equal(t, expOkCfg, cfg)
		},
	)

	t.Run(
		"Fail",
		func(t *testing.T) {
			fb := config.FallbackDecoder(errDec, errDec)(nil)
			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.Error(t, err)
			require.Empty(t, cfg)
		},
	)

	t.Run(
		"NilPrInChain",
		func(t *testing.T) {
			fb := config.FallbackDecoder(nil, okDecoder(reader()))(nil)
			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.NoError(t, err)
			require.Equal(t, expOkCfg, cfg)
		},
	)

	t.Run(
		"AllNils",
		func(t *testing.T) {
			fb := config.FallbackDecoder(nil, nil)(nil)
			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.Error(t, err)
			require.Empty(t, cfg)
		},
	)
}

type errDecoder struct{}

func (p *errDecoder) Decode(any) error {
	return errors.New("something went wrong")
}

func okDecoder(r io.Reader) config.Decoder {
	return json.NewDecoder(r)
}

func reader() *bytes.Reader {
	return bytes.NewReader([]byte(`{"foo": "bar"}`))
}
