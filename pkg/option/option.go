package option

type Option[T any] interface {
	IsSome() bool
	IsNone() bool

	Get() (T, bool)
	Unwrap() T
}

func Some[T any](v T) Option[T] {
	return some[T]{v}
}

func None[T any]() Option[T] {
	return none[T]{}
}

type some[T any] struct {
	v T
}

func (some[T]) IsSome() bool     { return true }
func (some[T]) IsNone() bool     { return false }
func (s some[T]) Get() (T, bool) { return s.v, true }
func (s some[T]) Unwrap() T      { return s.v }

type none[T any] struct{}

func (none[T]) IsSome() bool        { return false }
func (none[T]) IsNone() bool        { return true }
func (none[T]) Get() (t T, ok bool) { return t, false }
func (none[T]) Unwrap() (t T)       { return t }
