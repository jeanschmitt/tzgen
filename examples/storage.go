package main

import (
	"blockwatch.cc/tzgo/tezos"
	"fmt"
	"github.com/jeanschmitt/tzgen/examples/contracts"
)

// Storage example shows how to query the storage of a contract.
func (Examples) Storage() {
	client := NewClient(Ghostnet)
	fa2NFTAddress := tezos.MustParseAddress("KT1UvFD9hwgbLKbmTzUPmVRVsHG29ciQrcwE")

	fa2NFT, err := contracts.NewFA2NFT(ctx, fa2NFTAddress, client)
	fatalOnErr(err)

	storage, err := fa2NFT.Storage(ctx)
	fatalOnErr(err)

	// The returned struct uses go std, tzgo/tezos or tzgen/bind types.
	fmt.Println("Owner: ", storage.Owner)
	fmt.Println("OwnerCandidate: ", storage.OwnerCandidate)
	fmt.Println("Paused: ", storage.Paused)
	fmt.Println("Ledger: ", storage.Ledger)
}
