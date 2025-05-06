package config_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hikitani/go-config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCmdlineConfig struct {
	Verbose []bool `short:"v" long:"verbose"`
	Name    string `long:"name" required:"true"`
	Surname string `long:"surname" default:"foo"`
}

func TestCmdline(t *testing.T) {
	os.Args = []string{"program", "-vvv", "--name", "John"}

	cfg, err := config.New[testCmdlineConfig](
		config.FromCmdline(),
		config.WithDecoder(config.CmdlineDecoder),
	)
	require.NoError(t, err)
	assert.Equal(t, testCmdlineConfig{Verbose: []bool{true, true, true}, Name: "John", Surname: "foo"}, cfg)
}

func TestPrintHelp(t *testing.T) {
	const ExpectedUsage = `Usage:
  program [OPTIONS]

Application Options:
  -v, --verbose
      --name=
      --surname=

Help Options:
  -h, --help     Show this help message

`

	tmpfile, err := os.Create(path.Join(t.TempDir(), "test-stdout"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = tmpfile.Close() })
	os.Stdout = tmpfile

	defer func() {
		// os.Exit catching
		if v, ok := recover().(string); ok && strings.Contains(v, "os.Exit(0)") {
			actualUsage, err := os.ReadFile(tmpfile.Name())
			require.NoError(t, err)
			assert.Equal(t, ExpectedUsage, string(actualUsage))
		} else {
			t.Fatal(v)
		}
	}()

	os.Args = []string{"program", "--help"}
	_, _ = config.New[testCmdlineConfig](
		config.FromCmdline(),
		config.WithDecoder(config.CmdlineDecoder),
	)

	t.Fatal("For --help program must be terminated after creating config")
}
