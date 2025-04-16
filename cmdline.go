package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

const cmdSep = " "

var (
	_ ConfigProvider = &cmdlineProvider{}
	_ Decoder        = &cmdlineDecoder{}
)

type cmdlineProvider struct{}

// ProvideConfig implements ConfigProvider.
func (*cmdlineProvider) ProvideConfig() (io.Reader, error) {
	return strings.NewReader(strings.Join(os.Args[1:], cmdSep)), nil
}

// Returns config provider that provides config from cmdline arguments.
func FromCmdline() *cmdlineProvider {
	return &cmdlineProvider{}
}

type cmdlineDecoder struct {
	r io.Reader
}

// Decode implements Decoder.
func (d *cmdlineDecoder) Decode(v any) error {
	b, err := io.ReadAll(d.r)
	if err != nil {
		return fmt.Errorf("invalid reader: %w", err)
	}

	args := strings.Split(string(b), cmdSep)

	p := flags.NewParser(v, flags.HelpFlag|flags.PassDoubleDash)
	_, err = p.ParseArgs(args)
	if isErrHelp(err) {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(0)
	}
	if err != nil {
		return fmt.Errorf("parse args: %w", err)
	}

	return nil
}

// Returns cmdline decoder that parses cmd args to struct with tags.
// It uses under the hood [jessevdk/go-flags] library.
//
// [jessevdk/go-flags]: https://github.com/jessevdk/go-flags#example
func CmdlineDecoder(r io.Reader) *cmdlineDecoder {
	return &cmdlineDecoder{r: r}
}

func isErrHelp(err error) bool {
	if err == nil {
		return false
	}

	flagerr := &flags.Error{}
	if errors.As(err, &flagerr) {
		return flagerr.Type == flags.ErrHelp
	}

	return false
}
