package builder

import (
	"blockwatch.cc/tzgo/micheline"
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
)

func ToPrim(v interface{}) (micheline.Prim, error) {
	switch vv := v.(type) {
	case *big.Int:
		return intPrim(vv), nil
	case []byte:
		return bytesPrim(vv), nil
	case string:
		return stringPrim(vv), nil
	case bool:
		return boolPrim(vv), nil
	case Primer:
		return vv.ToPrim(), nil
	}

	return micheline.Prim{}, errors.Errorf("type not handled: %T", v)
}

type Primer interface {
	ToPrim() micheline.Prim
}

func intPrim(i *big.Int) micheline.Prim {
	return micheline.NewBig(i)
}

func bytesPrim(b []byte) micheline.Prim {
	return micheline.NewBytes(b)
}

func stringPrim(s string) micheline.Prim {
	return micheline.NewString(s)
}

func boolPrim(b bool) micheline.Prim {
	if b {
		return micheline.NewPrim(micheline.D_TRUE)
	}
	return micheline.NewPrim(micheline.D_FALSE)
}

func seqPrim(seq interface{}) micheline.Prim {
	seqVal := reflect.ValueOf(seq)
	if seqVal.Kind() != reflect.Slice && seqVal.Kind() != reflect.Array {
		panic("seq is neither a slice nor an array")
	}
	var prims []micheline.Prim
	for i := 0; i < seqVal.Len(); i++ {
		prim, err := ToPrim(seqVal.Index(i).Interface())
		if err != nil {
			panic(fmt.Sprintf("invalid primitive: %v", err))
		}
		prims = append(prims, prim)
	}
	return micheline.NewSeq(prims...)
}

func optionPrim(o interface{}) micheline.Prim {
	optVal := reflect.ValueOf(o)
	if optVal.Kind() != reflect.Ptr {
		panic("optionPrim expects a pointer arg")
	}
	if optVal.IsNil() {
		return micheline.NewPrim(micheline.D_NONE)
	}
	inner, err := ToPrim(optVal.Elem().Interface())
	if err != nil {
		panic("ToPrim failed")
	}
	prim := micheline.NewPrim(micheline.D_SOME)
	prim.Args = []micheline.Prim{inner}
	return prim
}
