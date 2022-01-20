package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"blockwatch.cc/tzgo/micheline"
)

func Parse(rawMicheline []byte, contractName string) (Contract, []*StructType, []*UnionType, []*OptionType, error) {
	c := &ParseContext{fmt: GoFormat}

	michelineBuf := new(bytes.Buffer)
	if err := json.Compact(michelineBuf, rawMicheline); err != nil {
		return Contract{}, nil, nil, nil, err
	}

	script := micheline.NewScript()
	err := json.Unmarshal(michelineBuf.Bytes(), &script)
	if err != nil {
		return Contract{}, nil, nil, nil, err
	}

	err = parseEntrypoints(c, script)
	if err != nil {
		return Contract{}, nil, nil, nil, err
	}

	return Contract{
		Name:        contractName,
		Micheline:   michelineBuf.String(),
		Entrypoints: c.entrypoints,
	}, c.structs, c.unions, c.options, nil
}

type ParseContext struct {
	entrypoints []*Entrypoint
	structs     []*StructType
	unions      []*UnionType
	options     []*OptionType
	counter     int

	fmt Formatter
}

func (c *ParseContext) inc() int {
	i := c.counter
	c.counter++
	return i
}

func (c *ParseContext) registerStruct(name string) *StructType {
	s := &StructType{Name: c.fmt.Struct(name)}
	c.structs = append(c.structs, s)
	return s
}

func (c *ParseContext) registerAutoStruct(refNames ...string) *StructType {
	name := fmt.Sprintf("%s%d", strings.Join(refNames, ""), c.inc())
	return c.registerStruct(name)
}

func (c *ParseContext) registerAutoUnion(refNames ...string) *UnionType {
	name := fmt.Sprintf("%sUnion%d", strings.Join(refNames, ""), c.inc())
	u := &UnionType{Name: c.fmt.Struct(name)}
	c.unions = append(c.unions, u)
	return u
}

func (c *ParseContext) registerAutoOption(refNames ...string) *OptionType {
	name := fmt.Sprintf("%sOption%d", strings.Join(refNames, ""), c.inc())
	o := &OptionType{Name: c.fmt.Struct(name)}
	c.options = append(c.options, o)
	return o
}

func (c *ParseContext) autoFieldName(refName string) string {
	return fmt.Sprintf("Field%d", c.inc())
}
