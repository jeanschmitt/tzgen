package parse

import (
	"blockwatch.cc/tzgo/micheline"
	"bytes"
	"encoding/json"
	"fmt"
	types "github.com/jeanschmitt/tzgen/pkg/types"
)

func Parse(rawMicheline []byte, contractName string) (*Contract, []*types.Struct, []*types.Union, error) {
	return NewParser(rawMicheline).Parse(contractName)
}

type Parser struct {
	script    *micheline.Script
	micheline []byte
	contract  *Contract

	structs map[[types.HashSize]byte]*types.Struct
	unions  []*types.Union
}

func NewParser(rawMicheline []byte) *Parser {
	return &Parser{
		micheline: rawMicheline,
	}
}

func (p *Parser) Parse(contractName string) (*Contract, []*types.Struct, []*types.Union, error) {
	var err error
	p.micheline, err = compactJson(p.micheline)
	if err != nil {
		return nil, nil, nil, err
	}

	// Init fields
	p.script = micheline.NewScript()
	p.contract = &Contract{
		Name:      contractName,
		Micheline: string(p.micheline),
	}
	p.structs = make(map[[types.HashSize]byte]*types.Struct)

	if err = json.Unmarshal(p.micheline, &p.script); err != nil {
		return nil, nil, nil, err
	}

	if err := p.entrypoints(); err != nil {
		return nil, nil, nil, err
	}

	var structs []*types.Struct
	i := 0
	for _, s := range p.structs {
		// Give a name to the structs
		s.Name = fmt.Sprintf("%sRecord%d", p.contract.Name, i)
		structs = append(structs, s)
		i++
	}

	for i, u := range p.unions {
		u.Name = fmt.Sprintf("%sUnion%d", p.contract.Name, i)
	}

	return p.contract, structs, p.unions, nil
}

func compactJson(src []byte) ([]byte, error) {
	dst := new(bytes.Buffer)
	if err := json.Compact(dst, src); err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
}

func (p *Parser) registerStruct(newStruct *types.Struct) *types.Struct {
	// Avoid struct duplicates using their hash
	var hash [types.HashSize]byte
	copy(hash[:], newStruct.Hash()[:types.HashSize])
	if found := p.lookupStruct(hash); found != nil {
		return found
	}
	p.structs[hash] = newStruct
	return newStruct
}

func (p *Parser) registerUnion(union *types.Union) *types.Union {
	// We don't avoid duplication for unions
	p.unions = append(p.unions, union)
	return union
}

func (p *Parser) lookupStruct(hash [types.HashSize]byte) *types.Struct {
	for h, s := range p.structs {
		if h == hash {
			return s
		}
	}
	return nil
}
