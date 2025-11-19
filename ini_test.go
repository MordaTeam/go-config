package config_test

import (
	"bytes"
	"testing"

	"github.com/MordaTeam/go-config"
	"github.com/stretchr/testify/require"
)

type iniCfg struct {
	Int    int     `ini:"INT"`
	Float  float64 `ini:"FLOAT64"`
	Bool   bool    `ini:"BOOL"`
	String string  `ini:"STRING"`
	Array  []int   `ini:"ARRAY"`
	Object Object  `ini:"OBJECT"`
}

type Object struct {
	FieldString string `ini:"FIELD_STRING"`
	FieldInt    int    `ini:"FIELD_INT"`
}

var data = []byte(`
INT = 42
FLOAT64 = 3.14159
BOOL = true
STRING = "Hello world"

ARRAY = 1,2,3,4,5

[OBJECT]
FIELD_STRING = "Nested value"
FIELD_INT = 100
`)

func TestIniDecoder(t *testing.T) {
	expCfg := iniCfg{
		Int:    42,
		Float:  3.14159,
		Bool:   true,
		String: "Hello world",
		Array:  []int{1, 2, 3, 4, 5},
		Object: Object{
			FieldString: "Nested value",
			FieldInt:    100,
		},
	}

	cfg, err := config.New[iniCfg](
		config.FromReader(bytes.NewReader(data)),
		config.WithDecoder(config.IniDecoder),
	)
	require.NoError(t, err)
	require.Equal(t, expCfg, cfg)
}
