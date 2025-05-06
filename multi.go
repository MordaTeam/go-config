package config

import (
	"errors"
	"fmt"
)

type multiConfigurator[T any] struct {
	configurators []configurator[T]
}

// OneOf builds config if at least once configurator created config successfully.
func (m *multiConfigurator[T]) OneOf() (cfg T, err error) {
	for i, c := range m.configurators {
		cerr := c(&cfg)
		if cerr == nil {
			return cfg, nil
		}

		err = errors.Join(err, fmt.Errorf("in pos %d: %w", i, cerr))
	}

	return cfg, fmt.Errorf("create config from configurators: %w", err)
}

// OneOfFill fills the provided config pointer with the result of OneOf method.
// Note that all fields of cfg will be overridden by new values.
func (m *multiConfigurator[T]) OneOfFill(cfg *T) (err error) {
	*cfg, err = m.OneOf()
	return
}

// AllOf builds config if all configurators created config successfully.
func (m *multiConfigurator[T]) AllOf() (cfg T, err error) {
	for i, c := range m.configurators {
		if cerr := c(&cfg); cerr != nil {
			err = errors.Join(err, fmt.Errorf("in pos %d: %w", i, cerr))
		}
	}

	if err != nil {
		var empty T
		return empty, fmt.Errorf("create config from configurators: %w", err)
	}

	return cfg, nil
}

// AllOfFill fills the provided config pointer with the result of AllOf method.
// Note that all fields of cfg will be overridden by new values.
func (m *multiConfigurator[T]) AllOfFill(cfg *T) (err error) {
	*cfg, err = m.AllOf()
	return
}

// Add adds configurator to build config.
func (m *multiConfigurator[T]) Add(provider ConfigProvider, opts ...ConfigOption) *multiConfigurator[T] {
	m.configurators = append(m.configurators, newConfigurator[T](provider, opts...))
	return m
}

// Multi creates multi configurator that aggregates different methods of creating config.
// Use method .Add to add configurator, then call .OneOf or .AllOf method to build config.
//
// Example
//
//	cfg, err := config.Multi[MyConfig]().
//		Add(config.FromReader(file)).
//		Add(config.FromEnv(), config.WithDecoder(config.EnvDecoder)).
//		Add(config.FromCmdline(), config.WithDecoder(config.CmdlineDecoder)).
//		AllOf()
func Multi[T any]() *multiConfigurator[T] {
	return &multiConfigurator[T]{}
}

type configurator[T any] func(cfg *T) error

func newConfigurator[T any](provider ConfigProvider, opts ...ConfigOption) configurator[T] {
	return func(cfg *T) error {
		return Fill(cfg, provider, opts...)
	}
}
