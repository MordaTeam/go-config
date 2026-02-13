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
	expOkCfg := fbConfig{Foo: "bar"}

	t.Run(
		"Normal_OkFirst",
		func(t *testing.T) {
			fb := config.FallbackDecoder(
				config.DecoderWrap(okDecoder),
				config.DecoderWrap(errDecoder),
			)(reader())

			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.NoError(t, err)
			require.Equal(t, expOkCfg, cfg)
		},
	)

	t.Run(
		"Normal_ErrFirst",
		func(t *testing.T) {
			fb := config.FallbackDecoder(
				config.DecoderWrap(errDecoder),
				config.DecoderWrap(okDecoder),
			)(reader())

			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.NoError(t, err)
			require.Equal(t, expOkCfg, cfg)
		},
	)

	t.Run(
		"Fail",
		func(t *testing.T) {
			fb := config.FallbackDecoder(
				config.DecoderWrap(errDecoder),
				config.DecoderWrap(errDecoder),
			)(reader())

			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.Error(t, err)
			require.Empty(t, cfg)
		},
	)

	t.Run(
		"NilPrInChain",
		func(t *testing.T) {
			fb := config.FallbackDecoder(
				nil,
				config.DecoderWrap(okDecoder),
			)(reader())

			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.NoError(t, err)
			require.Equal(t, expOkCfg, cfg)
		},
	)

	t.Run(
		"AllNils",
		func(t *testing.T) {
			fb := config.FallbackDecoder(
				nil,
				nil,
			)(reader())

			cfg := fbConfig{}
			err := fb.Decode(&cfg)
			require.Error(t, err)
			require.Empty(t, cfg)
		},
	)
}

type errDec struct{}

func (p *errDec) Decode(any) error {
	return errors.New("something went wrong")
}

func errDecoder(r io.Reader) *errDec {
	return &errDec{}
}

func okDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

func reader() *bytes.Reader {
	return bytes.NewReader([]byte(`{"foo": "bar"}`))
}
