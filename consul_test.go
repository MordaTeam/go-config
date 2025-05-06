package config_test

import (
	"context"
	"io"
	"testing"

	"github.com/MordaTeam/go-config"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/consul"
	"github.com/testcontainers/testcontainers-go/network"
)

func TestConsulProvider(t *testing.T) {
	ctx := context.Background()
	r := require.New(t)

	net, err := network.New(ctx)
	r.NoError(err)
	t.Cleanup(func() { _ = net.Remove(ctx) })

	consulContainer, err := consul.Run(ctx, "hashicorp/consul")
	r.NoError(err)
	t.Cleanup(func() { _ = consulContainer.Terminate(ctx) })

	r.NoError(consulContainer.Start(ctx))

	ip, err := consulContainer.Host(ctx)
	r.NoError(err)
	port, err := consulContainer.MappedPort(ctx, nat.Port("8500"))
	r.NoError(err)

	clientCfg := &api.Config{
		Address: ip + ":" + port.Port(),
	}
	client, err := api.NewClient(clientCfg)
	r.NoError(err)

	_, err = client.KV().Put(&api.KVPair{
		Key:   "foo/bar",
		Value: []byte(`{"foo": "bar"}`),
	}, &api.WriteOptions{})
	r.NoError(err)

	t.Run("DefaultClient", func(t *testing.T) {
		t.Setenv("CONSUL_HTTP_ADDR", clientCfg.Address)

		dataReader, err := config.
			FromConsul("/foo/bar").
			ProvideConfig()
		r.NoError(err)

		data, err := io.ReadAll(dataReader)
		r.NoError(err)

		r.Equal([]byte(`{"foo": "bar"}`), data)
	})

	t.Run("CustomClient", func(t *testing.T) {
		dataReader, err := config.
			FromConsul("/foo/bar", config.ConsulWithClient(client)).
			ProvideConfig()
		r.NoError(err)

		data, err := io.ReadAll(dataReader)
		r.NoError(err)

		r.Equal([]byte(`{"foo": "bar"}`), data)
	})
}
