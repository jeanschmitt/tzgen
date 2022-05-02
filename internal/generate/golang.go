package generate

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jeanschmitt/tzgen/pkg/ast/types"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"receiver":    receiver,
	"pascal":      strcase.ToCamel,
	"camel":       strcase.ToLowerCamel,
	"sub":         func(a, b int) int { return a - b },
	"type":        goType,
	"mkprim":      marshalPrimMethod,
	"pathFromIdx": pathFromIndex,
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
		return fmt.Sprintf("bind.Option[%s]", goType(t.Type))
	case *types.Union:
		return fmt.Sprintf("bind.Or[%s, %s]", goType(t.Left), goType(t.Right))
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
	//	return fmt.Sprintf("bind.Or[%s, %s]", goType(t.Left), goType(t.Right))
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

// pathFromIndex returns a path to a right-comb nested Pairs, from the index of a struct's field
// and the total number of fields.
func pathFromIndex(i, n int) string {
	if n == 1 {
		panic("pathFromIndex should not be called when a struct has 1 field")
	}
	// -
	// l r
	// l rl rr
	// l rl rrl rrr
	// l rl rrl rrrl rrrr
	if i == n-1 {
		return strings.TrimSuffix(strings.Repeat("r/", n-1), "/")
	}
	return strings.Repeat("r/", i) + "l"
}
