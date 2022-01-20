package codegen

import (
	"bytes"
	_ "embed"
	"go/format"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

// goTemplate contains the go contract template.
//go:embed golang.gotpl
var goTemplate []byte

type Data struct {
	Package  string
	Contract Contract
	Structs  []*StructType
	Unions   []*UnionType
	Options  []*OptionType
}

type Contract struct {
	Name        string
	Micheline   string
	Entrypoints []*Entrypoint
}

func (c Contract) Receiver() string {
	return strings.ToLower(string(c.Name[0]))
}

type Entrypoint struct {
	Name normalizedName
	Args []*Arg
}

type Arg struct {
	Name normalizedName
	Type Type
}

type Field = Arg

// region render

func (d *Data) Render(gofmt bool) ([]byte, error) {
	tpl, err := template.New("contract").Parse(string(goTemplate))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	// run gofmt
	out := buffer.Bytes()
	if gofmt {
		out, err = format.Source(out)
		if err != nil {
			return nil, errors.Wrap(err, "failed to format go code")
		}
	}

	return out, nil
}

// endregion
