package prim

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
)

func UnitV() string {
	return `{"prim":"Unit"}`
}

func IntV(i *big.Int) string {
	return `{"int":"` + i.String() + `"}`
}

func StringV(s string) string {
	return `{"string":"` + s + `"}`
}

func BytesV(b []byte) string {
	return `{"bytes":"` + hex.EncodeToString(b) + `"}`
}

func PairV(arg0, arg1 string) string {
	return `{"prim":"Pair","args":[` + arg0 + `,` + arg1 + `]}`
}

func ListV(args ...string) string {
	return "[" + strings.Join(args, ",") + "]"
}

func UnionV(arg string, branch UnionBranch) string {
	if branch == LeftBranch {
		return `{"prim":"Left","args":[` + arg + `]}`
	}
	return `{"prim":"Right","args":[` + arg + `]}`
}

func SomeV(arg string) string {
	return `{"prim":"Some","args":[` + arg + `]}`
}

func NoneV() string {
	return `"{prim":"None"}`
}

type Primer interface {
	ToPrim() string
}

type Int struct {
	i *big.Int
}

func NewInt(i *big.Int) Int {
	return Int{i: i}
}

func (i Int) ToPrim() string {
	return IntV(i.i)
}

type String struct {
	s string
}

func NewString(s string) String {
	return String{s: s}
}

func (s String) ToPrim() string {
	return StringV(s.s)
}

type Bytes struct {
	b []byte
}

func NewBytes(b []byte) Bytes {
	return Bytes{b: b}
}

func (b Bytes) ToPrim() string {
	return BytesV(b.b)
}

// region union

type UnionBranch bool

const (
	LeftBranch  UnionBranch = true
	RightBranch UnionBranch = false
)

func (u UnionBranch) IsLeft() bool {
	return u == LeftBranch
}

func (u UnionBranch) IsRight() bool {
	return u == RightBranch
}

// endregion

var ErrNone = errors.New("none")
