package prim

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
)

func unitV() string {
	return `{"prim":"Unit"}`
}

func intV(i *big.Int) string {
	return `{"int":"` + i.String() + `"}`
}

func stringV(s string) string {
	return `{"string":"` + s + `"}`
}

func bytesV(b []byte) string {
	return `{"bytes":"` + hex.EncodeToString(b) + `"}`
}

func pairV(arg0, arg1 string) string {
	return `{"prim":"Pair","args":[` + arg0 + `,` + arg1 + `]}`
}

func listV(args ...string) string {
	return "[" + strings.Join(args, ",") + "]"
}

func unionV(arg string, branch UnionBranch) string {
	if branch == LeftBranch {
		return `{"prim":"Left","args":[` + arg + `]}`
	}
	return `{"prim":"Right","args":[` + arg + `]}`
}

func someV(arg string) string {
	return `{"prim":"Some","args":[` + arg + `]}`
}

func noneV() string {
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
	return intV(i.i)
}

type String struct {
	s string
}

func NewString(s string) String {
	return String{s: s}
}

func (s String) ToPrim() string {
	return stringV(s.s)
}

type Bytes struct {
	b []byte
}

func NewBytes(b []byte) Bytes {
	return Bytes{b: b}
}

func (b Bytes) ToPrim() string {
	return bytesV(b.b)
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
