package main

import (
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/signer"
	"blockwatch.cc/tzgo/tezos"
	"fmt"
	"github.com/jeanschmitt/tzgen/examples/contracts"
)

// Call example shows how to call a contract's entry.
func (Examples) Call() {
	client := NewClient(Ghostnet)
	helloAddress := tezos.MustParseAddress("KT1FsGDhS7PUWn1Yff8r5nXE4MArZwD2XhEi")

	opts := rpc.DefaultOptions
	opts.Signer = signer.NewFromKey(AlicePrivateKey)

	hello, err := contracts.NewHello(ctx, helloAddress, client)
	fatalOnErr(err)

	fmt.Println("Calling greet...")

	tx, err := hello.Greet(ctx, &opts, "World!")
	fatalOnErr(err)

	fmt.Println("Operation hash: ", tx.Op.Hash)
	fmt.Printf("New storage: %q\n", must(hello.Storage(ctx)))
}
