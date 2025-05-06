package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/MordaTeam/go-config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEnvConfig struct {
	Secret string   `env:"SECRET"`
	Hosts  []string `env:"HOSTS"`
}

func TestEnv_NoProvider(t *testing.T) {
	expectedConfig := testEnvConfig{
		Secret: "the-most-secret-value",
		Hosts:  []string{"localhost", "foo", "bar"},
	}
	require.NoError(t, os.Setenv("SECRET", expectedConfig.Secret))
	require.NoError(t, os.Setenv("HOSTS", strings.Join(expectedConfig.Hosts, ",")))

	cfg, err := config.New[testEnvConfig](
		config.NoProvider, config.WithDecoder(config.EnvDecoder),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedConfig, cfg)
}

type testEnvConfigWithEnvProvider struct {
	Secret   string   `env:"SECRET"`
	Hosts    []string `env:"HOSTS"`
	Hostname string   `env:"HOSTNAME"`
}

func TestEnv_WithEnvProvider(t *testing.T) {
	hostname, err := os.Hostname()
	require.NoError(t, err)

	expectedConfig := testEnvConfigWithEnvProvider{
		Secret:   "the-most-secret-value",
		Hosts:    []string{"localhost", "foo", "bar"},
		Hostname: hostname,
	}
	require.NoError(t, os.Setenv("SECRET", expectedConfig.Secret))
	require.NoError(t, os.Setenv("HOSTS", strings.Join(expectedConfig.Hosts, ",")))

	cfg, err := config.New[testEnvConfigWithEnvProvider](
		config.FromEnv(), config.WithDecoder(config.EnvDecoder),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedConfig, cfg)
}

func TestEnv_WithEnvProvider_WithExpand(t *testing.T) {
	hostname, err := os.Hostname()
	require.NoError(t, err)

	expectedConfig := testEnvConfigWithEnvProvider{
		Secret:   "the-most-secret-value",
		Hosts:    []string{"localhost", "foo", "bar"},
		Hostname: hostname,
	}
	require.NoError(t, os.Setenv("VERY_SECRET", expectedConfig.Secret))
	require.NoError(t, os.Setenv("SECRET", "$VERY_SECRET"))
	require.NoError(t, os.Setenv("HOSTS", strings.Join(expectedConfig.Hosts, ",")))

	cfg, err := config.New[testEnvConfigWithEnvProvider](
		config.FromEnvWithOptions(config.EnvOptions{ExpandEnv: true}), config.WithDecoder(config.EnvDecoder),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedConfig, cfg)
}

type testConfigWithEnvProviderWithLowerCase struct {
	Secret   string   `env:"secret"`
	Hosts    []string `env:"hosts"`
	Hostname string   `env:"hostname"`
}

func TestEnv_WithEnvProvider_WithLowerCase(t *testing.T) {
	expectedConfig := testConfigWithEnvProviderWithLowerCase{
		Secret: "the-most-secret-value",
		Hosts:  []string{"localhost", "foo", "bar"},
	}
	require.NoError(t, os.Setenv("SECRET", expectedConfig.Secret))
	require.NoError(t, os.Setenv("HOSTS", strings.Join(expectedConfig.Hosts, ",")))

	cfg, err := config.New[testConfigWithEnvProviderWithLowerCase](
		config.FromEnvWithOptions(config.EnvOptions{KeyToLowerCase: true}), config.WithDecoder(config.EnvDecoder),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedConfig, cfg)
}
