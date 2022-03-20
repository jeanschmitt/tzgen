package tzgoext

import (
	"blockwatch.cc/tzgo/micheline"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
	"time"
)

func UnmarshalPrim(prim micheline.Prim, v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("v should be a non-nil pointer")
	}
	val = val.Elem()

	switch prim.Type {
	case micheline.PrimInt:
		return unmarshalInt(prim.Int, val)
	case micheline.PrimString:
		return unmarshalString(prim.String, val)
	case micheline.PrimBytes:
		return unmarshalBytes(prim.Bytes, val)
		// TODO: containers
	}
	return nil
}

func unmarshalInt(i *big.Int, v reflect.Value) error {
	typ := v.Type()
	if typ == tBigInt {
		v.Set(reflect.ValueOf(i))
	} else if typ == tTime {
		v.Set(reflect.ValueOf(time.Unix(i.Int64(), 0)))
	} else {
		return errors.Errorf("unexpected type for int prim: %T", v.Type())
	}
	return nil
}

func unmarshalString(str string, v reflect.Value) error {
	typ := v.Type()
	if typ == tString {
		v.SetString(str)
	} else {
		return errors.Errorf("unexpected type for string prim: %T", v.Type())
	}
	return nil
}

func unmarshalBytes(b []byte, v reflect.Value) error {
	typ := v.Type()
	if typ == tBytes {
		v.Set(reflect.ValueOf(b))
	} else {
		return errors.Errorf("unexpected type for string prim: %T", v.Type())
	}
	return nil
}

var (
	tBigInt = reflect.TypeOf((*big.Int)(nil))
	tString = reflect.TypeOf("")
	tTime   = reflect.TypeOf(time.Time{})
	tBytes  = reflect.TypeOf(([]byte)(nil))
)
