package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/MordaTeam/go-toolbox/options"
)

// ConfigProvider is an interface that provides configuration data.
type ConfigProvider interface {
	ProvideConfig() (io.Reader, error)
}

// Decoder is an interface that decodes configuration data into an object.
// The object must be pointer.
type Decoder interface {
	Decode(v any) error
}

type cfgOpts struct {
	newDec func(r io.Reader) Decoder
}

// WithDecoder is an option that overrides the default decoder. By default, it uses json.Decoder.
func WithDecoder[D Decoder](newDec func(r io.Reader) D) options.Option[cfgOpts] {
	return func(v *cfgOpts) error {
		v.newDec = func(r io.Reader) Decoder {
			return newDec(r)
		}
		return nil
	}
}

// Creates config T where provider provides data for decoding. By default, it uses json.Decoder.
// Use WithDecoder to override the decoder. If the reader implements the io.Closer interface, then
// it will be closed.
func New[T any](provider ConfigProvider, opts ...options.Option[cfgOpts]) (cfg T, err error) {
	r, err := provider.ProvideConfig()
	if err != nil {
		return cfg, fmt.Errorf("provide config: %w", err)
	}

	defer func() {
		r, ok := r.(io.Closer)
		if !ok {
			return
		}

		closeErr := r.Close()
		if closeErr != nil {
			closeErr = fmt.Errorf("close reader: %w", closeErr)
		}

		err = errors.Join(err, closeErr)
	}()

	cfgOpts := cfgOpts{
		newDec: func(r io.Reader) Decoder {
			return json.NewDecoder(r)
		},
	}

	for _, option := range opts {
		if option == nil {
			continue
		}
		if err := option(&cfgOpts); err != nil {
			return cfg, fmt.Errorf("apply option: %w", err)
		}
	}

	if err := cfgOpts.newDec(r).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("decode config: %w", err)
	}

	return cfg, nil
}

// Fill creates config and fills into cfg argument.
//
// Provider provides data for decoding. By default, it uses decoder json.Decoder.
// Use WithDecoder to override the decoder. If the reader implements the io.Closer interface, then
// it will be closed.
//
// For example we want to create config with default values and then fill config.
//
//	type Config struct {
//		Foo string `json:"foo"`
//	}
//
//	func main() {
//		cfg := Config{
//			Foo: "default"
//		}
//
//		err := config.Fill(&cfg, config.FromConsul"/bar/foo"))
//		//...
//	}
func Fill[T any](cfg *T, provider ConfigProvider, opts ...options.Option[cfgOpts]) error {
	newCfg, err := New[T](provider, opts...)
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}

	*cfg = mergeLeft(*cfg, newCfg)
	return nil
}
