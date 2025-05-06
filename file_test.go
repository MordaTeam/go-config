package config_test

import (
	"os"
	"path"
	"testing"

	"github.com/MordaTeam/go-config"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Foo string `json:"foo"`
}

func TestFile(t *testing.T) {
	filePath := path.Join(t.TempDir(), "config.json")
	file, err := os.Create(filePath)
	require.NoError(t, err)

	_, err = file.WriteString(`{"foo": "bar"}`)
	require.NoError(t, err)
	require.NoError(t, file.Close())

	cfg, err := config.New[testConfig](config.FromFile(filePath))
	require.NoError(t, err)
	require.Equal(t, "bar", cfg.Foo)
}
