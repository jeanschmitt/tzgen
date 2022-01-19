package codegen

import "github.com/pkg/errors"

func Generate(micheline []byte, packageName, contractName string) ([]byte, error) {
	contract, structs, unions, options, err := Parse(micheline, contractName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse micheline")
	}

	data := Data{
		Package:  packageName,
		Contract: contract,
		Structs:  structs,
		Unions:   unions,
		Options:  options,
	}

	res, err := data.Render(true)
	if err != nil {
		return nil, err
	}

	return res, nil
}
