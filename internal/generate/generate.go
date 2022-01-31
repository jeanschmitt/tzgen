package generate

import (
	"github.com/jeanschmitt/tzgen/internal/parse"
	"text/template"
)

func Generate(micheline []byte, contractName string, languageSpecific Metadata) ([]byte, error) {
	contract, structs, unions, err := parse.Parse(micheline, contractName)
	if err != nil {
		return nil, err
	}

	data := Data{
		Contract: contract,
		Structs:  structs,
		Unions:   unions,
		Metadata: languageSpecific,
	}

	err = languageSpecific.PreRender(&data)
	if err != nil {
		return nil, err
	}

	out, err := data.Render(languageSpecific.Language())
	if err != nil {
		return nil, err
	}

	return languageSpecific.PostRender(out)
}

type Metadata interface {
	Funcs() template.FuncMap
	PreRender(data *Data) error
	PostRender(in []byte) (out []byte, err error)
	Language() Language
}
