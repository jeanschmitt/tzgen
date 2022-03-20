package tzgoext

import (
	"blockwatch.cc/tzgo/micheline"
	"github.com/jeanschmitt/tzgen/pkg/either"
	"github.com/jeanschmitt/tzgen/pkg/option"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
	"time"
)

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

	// TODO: it seems that FoldPair doesn't work correctly
	return micheline.NewSeq(prims...).FoldPair(), nil
}

func MarshalPrim(v any) (micheline.Prim, error) {
	// Handle types that we can process with a type switch
	switch t := v.(type) {
	case micheline.PrimMarshaler:
		return t.MarshalPrim()
	case bool:
		return MarshalPrimBool(t), nil
	case string:
		return micheline.NewString(t), nil
	case *big.Int:
		return micheline.NewBig(t), nil
	case []byte:
		return micheline.NewBytes(t), nil
	case time.Time:
		return MarshalPrimTimestamp(t), nil
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Slice:

	}

	return micheline.Prim{}, errors.Errorf("type not handled: %T", v)
}

func MarshalPrimBool(b bool) micheline.Prim {
	if b {
		return micheline.NewCode(micheline.D_TRUE)
	}
	return micheline.NewCode(micheline.D_FALSE)
}

func MarshalPrimTimestamp(t time.Time, optimized ...bool) micheline.Prim {
	if len(optimized) == 1 && optimized[0] {
		return micheline.NewBig(big.NewInt(t.Unix()))
	}
	return micheline.NewString(t.Format(time.RFC3339))
}

func MarshalPrimOption[T any](o option.Option[T]) (micheline.Prim, error) {
	if o.IsSome() {
		inner, err := MarshalPrim(o.Unwrap())
		if err != nil {
			return micheline.Prim{}, err
		}
		return micheline.NewCode(micheline.D_SOME, inner), nil
	}
	return micheline.NewCode(micheline.D_NONE), nil
}

func MarshalPrimEither[L, R any](o either.Either[L, R], inner micheline.Prim) micheline.Prim {
	if o.IsLeft() {
		return micheline.NewCode(micheline.D_LEFT, inner)
	}
	return micheline.NewCode(micheline.D_RIGHT, inner)
}
