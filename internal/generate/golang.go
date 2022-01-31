package generate

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jeanschmitt/tzgen/pkg/types"
	"github.com/pkg/errors"
	"go/format"
	"strconv"
	"strings"
	"text/template"
)

type GoMetadata struct {
	Package  string
	Receiver string
}

func NewGoMetadata(pkg string) *GoMetadata {
	return &GoMetadata{
		Package: pkg,
	}
}

func (m *GoMetadata) Funcs() template.FuncMap {
	return map[string]interface{}{
		"struct":        strcase.ToCamel,
		"method":        strcase.ToCamel,
		"arg":           strcase.ToLowerCamel,
		"field":         strcase.ToCamel,
		"arglist":       m.argList(strcase.ToLowerCamel),
		"fieldlist":     m.argList(strcase.ToCamel),
		"type":          m.printableType,
		"receiver":      m.receiver,
		"isPtr":         m.isPtr,
		"builderMethod": m.builderMethod,
	}
}

type printableArg struct {
	Name          string
	Type          string
	BuilderMethod string
}

func (m *GoMetadata) argList(caseFunc func(string) string) func(args []types.Param) []printableArg {
	return func(args []types.Param) []printableArg {
		var printable []printableArg
		for i, arg := range args {
			if arg.Type.TypeName() == types.TypeUnit {
				continue
			}

			name := arg.Name
			if name == "" {
				name = arg.Type.TypeName() + strconv.Itoa(i)
			}
			printable = append(printable, printableArg{
				Name:          caseFunc(name),
				Type:          m.printableType(arg.Type),
				BuilderMethod: m.builderMethod(arg.Type),
			})
		}
		return printable
	}
}

func (m GoMetadata) Language() Language {
	return GoLanguage
}

func (m *GoMetadata) PreRender(data *Data) error {
	if len(data.Contract.Name) == 0 {
		return errors.New("empty contract name")
	}
	m.Receiver = m.receiver(data.Contract.Name)
	return nil
}

const skipGoFmt = false

func (m *GoMetadata) PostRender(in []byte) (out []byte, err error) {
	if skipGoFmt {
		return in, nil
	}
	out, err = format.Source(in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to format go code")
	}
	return out, nil
}

func (m *GoMetadata) printableType(typ types.Type) string {
	switch t := typ.(type) {
	case types.Nat, types.Int, types.Mutez:
		return "*big.Int"
	case types.Bytes, types.Key, types.KeyHash, types.Signature:
		return "[]byte"
	case types.String, types.Address:
		return "string"
	case types.Bool:
		return "bool"
	case *types.Option:
		underlying := m.printableType(t.Type)
		// Prepend with `*`, since we use pointers to handle options
		if underlying[0] != '*' {
			underlying = "*" + underlying
		}
		return underlying
	case *types.Struct:
		return fmt.Sprintf("*%s", t.Name)
	case *types.Union:
		return fmt.Sprintf("*%s", t.Name)
	}
	return "interface{}"
}

func (m *GoMetadata) builderMethod(typ types.Type) string {
	switch typ.(type) {
	case types.Nat, types.Int, types.Mutez, types.Timestamp:
		return "Int"
	case types.Bytes, types.Key, types.KeyHash, types.Signature:
		return "Bytes"
	case types.String, types.Address:
		return "String"
	case types.Bool:
		return "Bool"
	case *types.List, *types.Set:
		return "Seq"
	case *types.Map, *types.BigMap:
		return "Map"
	case *types.Struct, *types.Union:
		return "Primer"
	case *types.Option:
		return "Option"
	}
	return ""
}

func (m *GoMetadata) receiver(s string) string {
	return strings.ToLower(string(s[0]))
}

// ptr ensures that the type is a pointer.
func (m *GoMetadata) isPtr(typ string) bool {
	if typ == "" {
		return false
	}
	return typ[0] == '*'
}
