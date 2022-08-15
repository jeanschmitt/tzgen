package generate

import (
	"bytes"
	_ "embed"
	"github.com/jeanschmitt/tzgen/pkg/ast"
	"github.com/jeanschmitt/tzgen/pkg/ast/types"
	"github.com/pkg/errors"
	"go/format"
	"log"
	"text/template"
)

//go:embed contract.go.tmpl
var goTemplate string

type Data struct {
	Contract *ast.Contract
	Structs  []*types.Struct
	Unions   []*types.Union
	Address  string
	Package  string
}

func Render(data *Data) ([]byte, error) {
	tpl, err := template.New("contract").Funcs(funcMap).Parse(goTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	out, err := format.Source(buffer.Bytes())
	if err != nil {
		log.Println(err)
		out = buffer.Bytes()
	}

	return out, nil
}
