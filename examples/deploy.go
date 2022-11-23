package main

import (
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/signer"
	"fmt"
	"github.com/jeanschmitt/tzgen/examples/contracts"
)

// Deploy example shows how to deploy a contract from a binding generated with tzgen.
func (Examples) Deploy() {
	client := NewClient(Ghostnet)

	opts := rpc.DefaultOptions
	opts.Signer = signer.NewFromKey(AlicePrivateKey)

	fmt.Println("Deploying Hello contract...")

	tx, contract, err := contracts.DeployHello(ctx, &opts, client, "initial storage")
	fatalOnErr(err)

	fmt.Println("Contract address: ", contract.Address())
	fmt.Println("Operation hash: ", tx.Op.Hash)
	fmt.Printf("Storage: %q\n", must(contract.Storage(ctx)))
}
