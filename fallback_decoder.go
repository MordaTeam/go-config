package config

import (
	"errors"
	"fmt"
)

var _ Decoder = &fallbackDecoder{}

type fallbackDecoder struct {
	decs []Decoder
}

func FallbackDecoder(decs ...Decoder) *fallbackDecoder {
	return &fallbackDecoder{decs: decs}
}

func (dec *fallbackDecoder) Decode(v any) error {
	var resErr error
	for idx, d := range dec.decs {
		if d == nil {
			resErr = errors.Join(resErr, fmt.Errorf("nil decoder on pos %d", idx))
			continue
		}

		if err := d.Decode(v); err != nil {
			resErr = errors.Join(resErr, fmt.Errorf("decoding on pos %d: %w", idx, err))
			continue
		}

		return nil
	}

	return resErr
}
