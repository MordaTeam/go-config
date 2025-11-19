package config_test

import (
	"io"
	"os"
	"testing"
	"testing/iotest"

	"github.com/MordaTeam/go-config"
	"github.com/stretchr/testify/require"
)

type testMultiConfig struct {
	Foo  string `json:"foo" env:"FOO"`
	Bar  string `json:"bar"`
	Port int    `long:"port"`
}

func fileCfg(t testing.TB, cfg string) *os.File {
	f, err := os.CreateTemp(os.TempDir(), "cfg_*.json")
	require.NoError(t, err)
	t.Cleanup(func() { _ = f.Close() })

	nw, err := f.WriteString(cfg)
	require.Equal(t, len(cfg), nw)
	require.NoError(t, err)

	offset, err := f.Seek(0, io.SeekStart)
	require.Equal(t, int64(0), offset)
	require.NoError(t, err)

	return f
}

func TestMulti_AllOf(t *testing.T) {
	r := require.New(t)
	os.Setenv("FOO", "env_foo")

	os.Args = []string{"cli", "--port", "8080"}

	f := fileCfg(t, `{"foo": "file_foo", "bar": "file_bar"}`)

	// env -> file -> cmdline
	cfg, err := config.Multi[testMultiConfig]().
		Add(config.FromEnv(), config.WithDecoder(config.EnvDecoder)).
		Add(config.FromReader(f)).
		Add(config.FromCmdline(), config.WithDecoder(config.CmdlineDecoder)).
		AllOf()
	r.NoError(err)
	r.Equal(testMultiConfig{
		Foo:  "file_foo",
		Bar:  "file_bar",
		Port: 8080,
	}, cfg)

	f = fileCfg(t, `{"foo": "file_foo", "bar": "file_bar"}`)

	// file -> env -> cmdline
	cfg, err = config.Multi[testMultiConfig]().
		Add(config.FromReader(f)).
		Add(config.FromEnv(), config.WithDecoder(config.EnvDecoder)).
		Add(config.FromCmdline(), config.WithDecoder(config.CmdlineDecoder)).
		AllOf()
	r.NoError(err)
	r.Equal(testMultiConfig{
		Foo:  "env_foo",
		Bar:  "file_bar",
		Port: 8080,
	}, cfg)

	f = fileCfg(t, `{"foo": "file_foo", "bar": "file_bar"}`)
	r.NoError(f.Close())

	cfg, err = config.Multi[testMultiConfig]().
		Add(config.FromReader(f)).
		Add(config.FromEnv(), config.WithDecoder(config.EnvDecoder)).
		Add(config.FromCmdline(), config.WithDecoder(config.CmdlineDecoder)).
		AllOf()
	r.Error(err)
	r.Equal(testMultiConfig{}, cfg)

	_, err = config.Multi[testMultiConfig]().
		Add(config.FromReader(iotest.ErrReader(io.EOF))).
		Add(config.FromEnv(), config.WithDecoder(config.EnvDecoder)).
		Add(config.FromCmdline(), config.WithDecoder(config.CmdlineDecoder)).
		AllOf()
	r.ErrorIs(err, io.EOF)
}

func TestMulti_OneOf(t *testing.T) {
	r := require.New(t)
	os.Setenv("FOO", "env_foo")
	os.Setenv("BAR", "env_bar")

	os.Args = []string{"cli", "--port", "8080"}

	f := fileCfg(t, `{"foo": "file_foo", "bar": "file_bar"}`)

	cfg, err := config.Multi[testMultiConfig]().
		Add(config.FromEnv(), config.WithDecoder(config.EnvDecoder)).
		Add(config.FromReader(f)).
		Add(config.FromCmdline(), config.WithDecoder(config.CmdlineDecoder)).
		OneOf()
	r.NoError(err)
	r.Equal(testMultiConfig{
		Foo: "env_foo",
	}, cfg)

	cfg, err = config.Multi[testMultiConfig]().
		Add(config.FromReader(iotest.ErrReader(io.EOF))).
		Add(config.FromEnv(), config.WithDecoder(config.EnvDecoder)).
		Add(config.FromCmdline(), config.WithDecoder(config.CmdlineDecoder)).
		OneOf()
	r.NoError(err)
	r.Equal(testMultiConfig{
		Foo: "env_foo",
	}, cfg)

	cfg, err = config.Multi[testMultiConfig]().
		Add(config.FromReader(iotest.ErrReader(io.EOF))).
		OneOf()
	r.Error(err)
	r.Equal(testMultiConfig{}, cfg)
}
