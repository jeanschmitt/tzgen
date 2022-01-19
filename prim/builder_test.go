package prim

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressNatNat(t *testing.T) {
	want := `{
		"prim": "Pair",
		"args": [
			{ "string": "KT1XBtM8Ww34LQ2UfXRqAFrpZfWKbz7UR4Yd" },
			{
				"prim": "Pair",
				"args": [
					{ "int": "1" },
					{ "int": "2" }
				]
			}
		]
	}`
	got := Builder().Address("KT1XBtM8Ww34LQ2UfXRqAFrpZfWKbz7UR4Yd").Int(big.NewInt(1)).Int(big.NewInt(2)).Finish()
	require.JSONEq(t, want, got)
}
