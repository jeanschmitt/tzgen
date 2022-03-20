package tzgoext

import (
	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/micheline"
)

type CallArgs struct {
	contract.TxArgs
	Entrypoint string
	Value      micheline.Prim
}

func (c CallArgs) Parameters() *micheline.Parameters {
	return &micheline.Parameters{
		Entrypoint: c.Entrypoint,
		Value:      c.Value,
	}
}

func (c CallArgs) Encode() *codec.Transaction {
	return &codec.Transaction{
		Manager: codec.Manager{
			Source: c.Source,
		},
		Destination: c.Destination,
		Parameters:  c.Parameters(),
	}
}
