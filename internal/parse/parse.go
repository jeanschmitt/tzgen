package parse

import (
	"blockwatch.cc/tzgo/micheline"
	"encoding/json"
	"fmt"
	"github.com/jeanschmitt/tzgen/pkg/ast"
	"github.com/jeanschmitt/tzgen/pkg/ast/types"
	"github.com/pkg/errors"
	"io"
	"strconv"
)

func Parse(codeReader io.Reader, name string) (*ast.Contract, []*types.Struct, []*types.Union, error) {
	code, err := io.ReadAll(codeReader)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to read code")
	}

	return newParser(code).parse(name)
}

type parser struct {
	script   *micheline.Script
	raw      []byte
	contract *ast.Contract

	counter int
	structs []structWithHash
	unions  []*types.Union
}

type structWithHash struct {
	*types.Struct
	hash [types.HashSize]byte
}

func newParser(raw []byte) *parser {
	return &parser{
		script:   micheline.NewScript(),
		raw:      raw,
		contract: new(ast.Contract),
	}
}

func (p *parser) parse(name string) (*ast.Contract, []*types.Struct, []*types.Union, error) {
	err := json.Unmarshal(p.raw, &p.script)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to unmarshal micheline code")
	}

	// Remove storage
	p.script.Storage = micheline.Prim{}
	p.raw, err = json.Marshal(p.script)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to re-marshall script")
	}

	p.contract.Name = name
	p.contract.Micheline = string(p.raw)

	if err = p.parseEntrypoints(); err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to parse entrypoints")
	}

	if p.contract.Storage, err = p.parseType(p.script.StorageType().TypedefPtr("Storage")); err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to parse storage")
	}

	return p.contract, p.nameStructs(), p.unions, nil
}

func (p *parser) nameStructs() []*types.Struct {
	var structs []*types.Struct
	for i, s := range p.structs {
		// Give a name to the structs
		if s.Name == "" {
			s.Name = fmt.Sprintf("%s_record_%d", p.contract.Name, i)
		} else {
			newName := fmt.Sprintf("%s_%s", p.contract.Name, s.Name)
			if p.structNameExists(newName) {
				newName = newName + strconv.Itoa(i)
			}
			s.Name = newName
		}
		structs = append(structs, s.Struct)
	}
	return structs
}

func (p *parser) inc() int {
	c := p.counter
	p.counter++
	return c
}

func (p *parser) registerStruct(newStruct *types.Struct) *types.Struct {
	// Avoid struct duplicates using their hash
	var hash [types.HashSize]byte
	copy(hash[:], newStruct.Hash()[:types.HashSize])
	if found := p.lookupStruct(hash); found != nil {
		return found
	}
	p.structs = append(p.structs, structWithHash{Struct: newStruct, hash: hash})
	return newStruct
}

func (p *parser) registerUnion(union *types.Union) *types.Union {
	// We don't avoid duplication for unions
	p.unions = append(p.unions, union)
	return union
}

func (p *parser) lookupStruct(hash [types.HashSize]byte) *types.Struct {
	for _, s := range p.structs {
		if s.hash == hash {
			return s.Struct
		}
	}
	return nil
}

func (p *parser) structNameExists(name string) bool {
	for _, s := range p.structs {
		if s.Name == name {
			return true
		}
	}
	return false
}
