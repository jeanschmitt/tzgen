package main

import (
	"blockwatch.cc/tzgo/tezos"
	"fmt"
	"github.com/jeanschmitt/tzgen/examples/contracts"
	"math/big"
)

// Bigmap example shows how to get bigmap handles from storage, and how to use them.
func (Examples) Bigmap() {
	client := NewClient(Ghostnet)
	fa2NFTAddress := tezos.MustParseAddress("KT1UvFD9hwgbLKbmTzUPmVRVsHG29ciQrcwE")

	fa2NFT, err := contracts.NewFA2NFT(ctx, fa2NFTAddress, client)
	fatalOnErr(err)

	storage, err := fa2NFT.Storage(ctx)
	fatalOnErr(err)

	fmt.Println("Ledger bigmap id: ", storage.Ledger.ID())

	val1, err := storage.Ledger.Get(ctx, big.NewInt(1))
	fatalOnErr(err)
	fmt.Println("Value for key 1: ", val1)

	_, err = storage.Ledger.Get(ctx, big.NewInt(-1))
	fmt.Println("Key -1 does not exist: ", err)
}
