package config

import (
	"fmt"
	"io"

	"github.com/go-ini/ini"
)

var _ Decoder = &iniDecoder{}

type iniDecoder struct {
	reader io.Reader
}

// Decode implements Decoder.
func (dec *iniDecoder) Decode(v any) error {
	if err := ini.MapTo(v, dec.reader); err != nil {
		return fmt.Errorf("mapping ini: %w", err)
	}
	return nil
}

func IniDecoder(r io.Reader) *iniDecoder {
	return &iniDecoder{reader: r}
}
