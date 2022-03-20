package ast

import (
	"blockwatch.cc/tzgo/micheline"
	"github.com/jeanschmitt/tzgen/pkg/ast/types"
)

type Entrypoint struct {
	Name   string
	Raw    *micheline.Entrypoint
	Params []*types.Param
}

type Getter struct {
	Entrypoint
	ReturnType types.Type
}
