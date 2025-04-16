package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	consul "github.com/hashicorp/consul/api"
)

var _ ConfigProvider = &consulProvider{}

type ConsulOption func(*consulOpts) error

type consulOpts struct {
	client  *consul.Client
	qryOpts *consul.QueryOptions
}

type consulProvider struct {
	client   *consul.Client
	qryOpts  *consul.QueryOptions
	cfgPath  string
	funcOpts []ConsulOption
}

func (c *consulProvider) lazyInit() error {
	if c.client != nil {
		return nil
	}

	var consulOpts consulOpts
	for _, option := range c.funcOpts {
		if option == nil {
			continue
		}

		if err := option(&consulOpts); err != nil {
			return fmt.Errorf("apply option: %w", err)
		}

	}

	if consulOpts.client == nil {
		var err error
		consulOpts.client, err = consul.NewClient(consul.DefaultConfig())
		if err != nil {
			return fmt.Errorf("create consul client: %w", err)
		}

	}

	c.client = consulOpts.client

	return nil
}

// ProvideConfig implements ConfigProvider
func (c *consulProvider) ProvideConfig() (io.Reader, error) {
	if err := c.lazyInit(); err != nil {
		return nil, fmt.Errorf("init lazy: %w", err)
	}

	qry := c.qryOpts
	if qry == nil {
		qry = &consul.QueryOptions{}
	}

	kv, _, err := c.client.KV().Get(c.cfgPath, qry)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, fmt.Errorf("consul get kv: key '%s' doesn't exist", c.cfgPath)
	}

	return bytes.NewBuffer(kv.Value), nil
}

// Overrides consul client.
func ConsulWithClient(client *consul.Client) ConsulOption {
	return func(v *consulOpts) error {
		if client == nil {
			return errors.New("got nil consul client")
		}

		v.client = client
		return nil
	}
}

// Defines the query options that will be used when lookuping for the config.
func ConsulWithQueryOptions(qry *consul.QueryOptions) ConsulOption {
	return func(v *consulOpts) error {
		if qry == nil {
			return errors.New("got nil consul query options")
		}

		v.qryOpts = qry
		return nil
	}
}

// Returns config provider that provides config from consul kv.
//
// If client wasn't passed with options, it's created with default config.
//
// By default, config will use CONSUL_HTTP_ADDR env as HTTP address.
// If it's empty, localhost will be chosen.
func FromConsul(cfgPath string, opts ...ConsulOption) *consulProvider {
	return &consulProvider{
		cfgPath:  cfgPath,
		funcOpts: opts,
	}
}
