package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/caarlos0/env/v9"
)

var (
	_ Decoder        = &envDecoder{}
	_ ConfigProvider = &envProvider{}
)

type EnvOptions struct {
	// Enable replacing ${var} or $var in the string according to the values of the current environment variables.
	// By default, false.
	ExpandEnv bool

	// Env keys will be converted to lower case (FOO -> foo).
	// By default, false.
	KeyToLowerCase bool
}

type envProvider struct {
	expand         bool
	keyToLowerCase bool
}

// ProvideConfig implements ConfigProvider.
func (e *envProvider) ProvideConfig() (io.Reader, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname: %w", err)
	}

	env := os.Environ()
	mapEnv := map[string]string{
		"HOSTNAME": hostname,
	}
	for _, entry := range env {
		key, val, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}

		if e.expand {
			val = os.ExpandEnv(val)
		}

		if e.keyToLowerCase {
			key = strings.ToLower(key)
		}

		mapEnv[key] = val
	}

	var r bytes.Buffer
	if err := json.NewEncoder(&r).Encode(mapEnv); err != nil {
		return nil, fmt.Errorf("encode env to json: %w", err)
	}

	return &r, nil
}

// FromEnvWithOptions returns provider that provides env variables in json form.
//
// Example:
//
//	// Env
//	FOO=bar
//	BUZ=foo
//	// converted to
//	{"FOO": "bar", "BUZ": "foo"}
func FromEnvWithOptions(opts EnvOptions) *envProvider {
	return &envProvider{
		expand:         opts.ExpandEnv,
		keyToLowerCase: opts.KeyToLowerCase,
	}
}

// FromEnv returns provider that provides env variables in json form.
//
// Example:
//
//	// Env
//	FOO=bar
//	BUZ=foo
//	// converted to
//	{"FOO": "bar", "BUZ": "foo"}
func FromEnv() *envProvider {
	return &envProvider{}
}

type envDecoder struct {
	mapEnv map[string]string
}

// Decode implements Decoder.
func (e *envDecoder) Decode(v any) error {
	if err := env.ParseWithOptions(v, env.Options{
		Environment: e.mapEnv,
	}); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	return nil
}

// Returns env decoder that parses envs to struct with tags.
// Provider should return json config or nothing.
// JSON config will be used as additional environment.
// From JSON config will be read to map[string]string, so variables must be string (other will be ignored).
// Use this decoder with provider NoProvider if you want variables only from environment.
// Recommended use with EnvProvider, because it provides variable $HOSTNAME.
// It uses under the hood [caarlos0/env/v9] library.
//
// [caarlos0/env/v9]: https://github.com/caarlos0/env
func EnvDecoder(r io.Reader) *envDecoder {
	var (
		dec     envDecoder
		jsonCfg map[string]any
	)
	if err := json.NewDecoder(r).Decode(&jsonCfg); err != nil {
		return &dec
	}

	mapEnv := map[string]string{}
	for k, v := range jsonCfg {
		if v, ok := v.(string); ok {
			mapEnv[k] = v
		}
	}

	dec.mapEnv = mapEnv
	return &dec
}
