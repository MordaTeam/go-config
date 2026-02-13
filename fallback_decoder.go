package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var _ Decoder = &fallbackDecoder{}

type fallbackDecoder struct {
	decCtrs []func(io.Reader) Decoder
	reader  io.Reader
}

// Creates new fallback line of decoders.
// Call for Decode() returns first successfull result of Decode() call of internal decoder.
// If all internal decoders fail, return resulting error.
func FallbackDecoder(decCtrs ...func(io.Reader) Decoder) func(r io.Reader) *fallbackDecoder {
	return func(r io.Reader) *fallbackDecoder {
		return &fallbackDecoder{
			decCtrs: decCtrs,
			reader:  r,
		}
	}
}

func (dec *fallbackDecoder) Decode(v any) error {
	data, err := io.ReadAll(dec.reader)
	if err != nil {
		return fmt.Errorf("reading data from reader: %w", err)
	}

	var resErr error
	for idx, dCtr := range dec.decCtrs {
		if dCtr == nil {
			resErr = errors.Join(resErr, fmt.Errorf("nil decoder constructor on pos %d", idx))
			continue
		}

		d := dCtr(bytes.NewReader(data))

		if err := d.Decode(v); err != nil {
			resErr = errors.Join(resErr, fmt.Errorf("decoding on pos %d: %w", idx, err))
			continue
		}

		return nil
	}

	return resErr
}

// Use this function to pass decoders to FallbackDecoder.
func DecoderWrap[D Decoder](src func(io.Reader) D) func(io.Reader) Decoder {
	return func(r io.Reader) Decoder {
		return src(r)
	}
}
