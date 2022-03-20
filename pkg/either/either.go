package either

type Either[L, R any] interface {
	IsLeft() bool
	IsRight() bool

	Left() (L, bool)
	Right() (R, bool)
}

func Left[L, R any](l L) Either[L, R] {
	return left[L, R]{v: l}
}

func Right[L, R any](r R) Either[L, R] {
	return right[L, R]{v: r}
}

type left[L, R any] struct {
	v L
}

func (left[L, R]) IsLeft() bool          { return true }
func (left[L, R]) IsRight() bool         { return false }
func (l left[L, R]) Left() (L, bool)     { return l.v, true }
func (left[L, R]) Right() (b R, ok bool) { return b, false }

type right[L, R any] struct {
	v R
}

func (right[L, R]) IsLeft() bool         { return false }
func (right[L, R]) IsRight() bool        { return true }
func (right[L, R]) Left() (a L, ok bool) { return a, false }
func (r right[L, R]) Right() (R, bool)   { return r.v, true }
