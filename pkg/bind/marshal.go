package bind

import (
	"blockwatch.cc/tzgo/micheline"
	"blockwatch.cc/tzgo/tezos"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
	"time"
)

// MarshalPrim marshals v into a Prim by using reflection.
func MarshalPrim(v any) (micheline.Prim, error) {
	// Handle types that we can process with a type switch
	switch t := v.(type) {
	case micheline.PrimMarshaler:
		return t.MarshalPrim()
	case *big.Int:
		return micheline.NewBig(t), nil
	case string:
		return micheline.NewString(t), nil
	case bool:
		if t {
			return micheline.NewCode(micheline.D_TRUE), nil
		}
		return micheline.NewCode(micheline.D_FALSE), nil
	case []byte:
		return micheline.NewBytes(t), nil
	case time.Time:
		return micheline.NewString(t.Format(time.RFC3339)), nil
	case tezos.Address:
		return micheline.NewString(t.String()), nil
	case tezos.Key:
		return micheline.NewString(t.String()), nil
	case tezos.Signature:
		return micheline.NewString(t.String()), nil
	case tezos.ChainIdHash:
		return micheline.NewString(t.String()), nil
	}

	// Container types
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Slice:
		n := val.Len()
		prims := make([]micheline.Prim, 0, n)
		for i := 0; i < n; i++ {
			prim, err := MarshalPrim(val.Index(i).Interface())
			if err != nil {
				return micheline.Prim{}, err
			}
			prims = append(prims, prim)
		}
		return micheline.NewSeq(prims...), nil
	}

	return micheline.Prim{}, errors.Errorf("type not handled: %T", v)
}

// MarshalParams marshals the provided params into a folded Prim.
func MarshalParams(params ...any) (micheline.Prim, error) {
	if len(params) == 0 {
		return micheline.NewPrim(micheline.D_UNIT), nil
	}

	prims := make([]micheline.Prim, 0, len(params))
	for _, p := range params {
		prim, err := MarshalPrim(p)
		if err != nil {
			return micheline.Prim{}, err
		}
		prims = append(prims, prim)
	}

	return foldRightComb(prims...), nil
}

// foldRightComb folds a list of prims into nested Pairs, by following the right-comb convention.
func foldRightComb(prims ...micheline.Prim) micheline.Prim {
	n := len(prims)
	switch n {
	case 0:
		return micheline.NewPrim(micheline.D_UNIT)
	case 1:
		return prims[0]
	default:
		return foldRightComb(append(prims[:n-2], micheline.NewPair(prims[n-2], prims[n-1]))...)
	}
}
