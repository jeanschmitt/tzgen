package generate

import (
	"bytes"
	_ "embed"
	"github.com/jeanschmitt/tzgen/internal/parse"
	"github.com/jeanschmitt/tzgen/pkg/types"
	"github.com/pkg/errors"
	"text/template"
)

// Data is the struct provided to the template.
type Data struct {
	Contract *parse.Contract
	Structs  []*types.Struct
	Unions   []*types.Union
	// Metadata contains language-specific data.
	// For example, it contains the package name in Go.
	Metadata Metadata
}

// Templates
var (
	//go:embed contract.go.tmpl
	goTemplate string
)

func (d *Data) Render(language Language) ([]byte, error) {
	tpl, err := template.New(language.String()).Funcs(d.Metadata.Funcs()).Parse(goTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	return buffer.Bytes(), nil
}

type Language string

const (
	GoLanguage Language = "go"
)

func (l Language) String() string {
	return string(l)
}
