package bind

import (
	"blockwatch.cc/tzgo/micheline"
	"github.com/pkg/errors"
)

// Option is a type that can either contain a value T, or be None.
type Option[T any] struct {
	v      T
	isSome bool
}

// Some returns a Some option with v as a value.
func Some[T any](v T) Option[T] {
	return Option[T]{v: v, isSome: true}
}

// None returns a None option for type T.
func None[T any]() Option[T] {
	return Option[T]{isSome: false}
}

// Get returns the inner value of the Option and a boolean
// indicating if the Option is Some.
//
// If it is none, the returned value is the default value for T.
func (o Option[T]) Get() (v T, isSome bool) {
	return o.v, o.isSome
}

// Unwrap returns the inner value of the Option, expecting
// that it is Some.
//
// Panics if the option is None.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		panic("Unwrap() called on a `None` Option")
	}
	return o.v
}

// UnwrapOr returns the inner value of the Option if it is Some,
// or the provided default value if it is None.
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.IsNone() {
		return defaultValue
	}
	return o.v
}

// UnwrapOrZero returns the inner value of the Option if it is Some,
// or T's zero value if it is None.
func (o Option[T]) UnwrapOrZero() T {
	// o.v == zero value if o is none
	return o.v
}

func (o Option[T]) IsSome() bool {
	return o.isSome
}

func (o Option[T]) IsNone() bool {
	return !o.isSome
}

func (o Option[T]) MarshalPrim(optimized bool) (micheline.Prim, error) {
	if o.isSome {
		inner, err := MarshalPrim(o.v, optimized)
		if err != nil {
			return micheline.Prim{}, err
		}
		return micheline.NewCode(micheline.D_SOME, inner), nil
	}
	return micheline.NewCode(micheline.D_NONE), nil
}

func (o *Option[T]) UnmarshalPrim(prim micheline.Prim) error {
	switch prim.OpCode {
	case micheline.D_SOME:
		if len(prim.Args) != 1 {
			return errors.New("prim Some should have 1 arg")
		}
		o.isSome = true
		return UnmarshalPrim(prim.Args[0], &o.v)
	case micheline.D_NONE:
		*o = None[T]()
		return nil
	default:
		return errors.Errorf("unexpected opCode when unmarshalling Option: %s", prim.OpCode)
	}
}

func (o Option[T]) keyHash() hashType {
	if v, ok := o.Get(); ok {
		return hashFunc(zero[T]())(v)
	} else {
		return hashType{0}
	}
}
