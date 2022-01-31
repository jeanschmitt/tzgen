package parse

import (
	types "github.com/jeanschmitt/tzgen/pkg/types"
)

type Contract struct {
	// Name of the contract.
	// It is not included in the contract's Michelson, so it should be provided
	// when running the tzgen.
	//
	// TODO: it is possible to find it for contracts implementing TZIP-16.
	// https://gitlab.com/tezos/tzip/-/blob/master/proposals/tzip-16/tzip-16.md
	Name string
	// Micheline contains the raw micheline code, used for generating the binding.
	Micheline   string
	Entrypoints []Entrypoint
	Bigmaps     []Bigmap
}

type Entrypoint struct {
	// Original Name
	// For single entrypoint contracts, the name is not provided in Michelson.
	// In this case, the entrypoint will be called `default`.
	Name string
	Args []types.Param
	id   int
}

type Bigmap struct {
	// Bigmap ID
	ID int64
	// Original Name of the bigmap, if provided.
	Name  string
	Types types.BigMap
}
