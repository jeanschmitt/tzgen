package main

import (
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/signer"
	"blockwatch.cc/tzgo/tezos"
	"fmt"
	"github.com/jeanschmitt/tzgen/examples/contracts"
)

// Payable example shows how to call a contract entry with an Amount value different from 0.
func (Examples) Payable() {
	client := NewClient(Ghostnet)
	payableAddress := tezos.MustParseAddress("KT1FPA7vN4cBk24df7VxUu9DstRcc7am3qnf")

	opts := rpc.DefaultOptions
	opts.Signer = signer.NewFromKey(AlicePrivateKey)

	payable, err := contracts.NewPayable(ctx, payableAddress, client)
	fatalOnErr(err)

	fmt.Println("Calling send_tz with an amount of 42µtz...")

	params, err := payable.Builder().SendTz()
	fatalOnErr(err)

	tx, err := payable.Call(ctx, &contract.TxArgs{Params: params, Amount: 42}, &opts)
	fatalOnErr(err)

	fmt.Println("Operation hash: ", tx.Op.Hash)
	fmt.Printf("Received: %vµtz\n", must(payable.Storage(ctx)))
}
