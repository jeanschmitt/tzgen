package main

import (
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"context"
	"fmt"
	"github.com/jeanschmitt/tzgen/examples/contracts"
	"os"
)

//go:generate go run ../main.go --address KT1CiYNu9iJknnL31TXBWHCqRdFRh7jPWdzg --name Simple --pkg contracts -o ./contracts/Simple.go

var examples = map[string]func(){
	"tzip4": TZIP4View,
}

func TZIP4View() {
	simple := simpleContract(rpcClient())

	counter, err := simple.GetCounter(context.Background())
	handleErr(err)

	fmt.Printf("Current counter value: %v\n", counter)
}

func main() {
	if len(os.Args) > 1 {
		for _, name := range os.Args[1:] {
			ex, ok := examples[name]
			if !ok {
				fmt.Printf("Example %s does not exist\n", name)
				continue
			}
			fmt.Printf("\nRunning %s example...\n", name)
			ex()
		}
	} else {
		for name, ex := range examples {
			fmt.Printf("\nRunning %s example...\n", name)
			ex()
		}
	}
}

func rpcClient() *rpc.Client {
	tzClient, err := rpc.NewClient("https://hangzhounet.smartpy.io", nil)
	tzClient.ChainId = tezos.Hangzhounet2
	handleErr(err)
	return tzClient
}

func simpleContract(tzClient *rpc.Client) *contracts.Simple {
	return contracts.NewSimple(tezos.MustParseAddress("KT1CiYNu9iJknnL31TXBWHCqRdFRh7jPWdzg"), tzClient)
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("An error occured: %v\n", err)
		os.Exit(1)
	}
}
