package generate

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jeanschmitt/tzgen/pkg/ast/types"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"receiver": receiver,
	"pascal":   strcase.ToCamel,
	"camel":    strcase.ToLowerCamel,
	"sub":      func(a, b int) int { return a - b },
	"type":     goType,
	"mkprim":   marshalPrimMethod,
}

func receiver(typeName string) string {
	if typeName == "" {
		return "_r"
	}
	return "_" + strings.ToLower(string(typeName[0]))
}

func goType(typ types.Type) string {
	switch t := typ.(type) {
	case types.Nat,
		types.Int,
		types.Mutez:
		return "*big.Int"
	case types.String:
		return "string"
	case types.Bool:
		return "bool"
	case types.Bytes,
		types.KeyHash:
		return "[]byte"
	case types.Timestamp:
		return "time.Time"
	case types.Address:
		return "tezos.Address"
	case types.Key:
		return "tezos.Key"
	case types.Signature:
		return "tezos.Signature"
	case types.ChainID:
		return "tezos.ChainIdHash"

	case *types.Option:
		return fmt.Sprintf("option.Option[%s]", goType(t.Type))
	case *types.Union:
		return fmt.Sprintf("either.Either[%s, %s]", goType(t.Left), goType(t.Right))
	case *types.List:
		return "[]" + goType(t.Type)
	case *types.Set:
		return "[]" + goType(t.Type)
	case *types.Map:
		//return fmt.Sprintf("*hashmap.Map[%s, %s]", goType(t.Key), goType(t.Value))
	case *types.Struct:
		return "*" + strcase.ToCamel(t.Name)

		// case types.Unit:
		// case types.Operation:
		// case types.Contract:
	}

	return "any"
}

func marshalPrimMethod(typ types.Type) string {
	switch t := typ.(type) {
	case types.Nat,
		types.Int,
		types.Mutez:
		return "micheline.NewBig(%s)"
	case types.String:
		return "micheline.NewString(%s)"
	case types.Bool:
		return "tzgoext.MarshalPrimBool(%s)"
	case types.Bytes,
		types.KeyHash:
		return "micheline.NewBytes(%s)"
	case types.Timestamp:
		return "tzgoext.MarshalPrimTimestamp(%s)"
	case types.Address:
		return "micheline.NewString(%s.String())"
	case types.Key,
		types.Signature,
		types.ChainID:
		return "micheline.NewBytes(%s.Bytes())"

	//case *types.Option:
	//	return "tzgoext.MarshalPrimOption(%s"
	//case *types.Union:
	//	return fmt.Sprintf("either.Either[%s, %s]", goType(t.Left), goType(t.Right))
	case *types.List:
		return "tzgoext.MarshalPrimSeq[" + goType(t.Type) + "](%s, tzgoext.MarshalAny)"
	//case *types.Set:
	//	return "[]" + goType(t.Type)
	//case *types.Map:
	//	return fmt.Sprintf("tzgoext.Map[%s, %s]", goType(t.Key), goType(t.Value))
	case *types.Struct:
		return "%s.Prim()"

	default:
		// case *types.Unit:
		// case *types.Operation:
		// case *types.Contract:
	}

	return "micheline.Prim{}/* %s */"
}
