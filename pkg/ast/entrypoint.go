package ast

import (
	"blockwatch.cc/tzgo/micheline"
	"github.com/jeanschmitt/tzgen/pkg/ast/types"
)

type Entrypoint struct {
	Name   string
	Raw    *micheline.Entrypoint `json:"-"`
	Params []*types.Param
}

// Getter is a read-only Entrypoint, with a return value.
//
// It is implemented with TZIP-4.
type Getter struct {
	Entrypoint
	ReturnType types.Type
}
