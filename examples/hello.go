package main

import (
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"context"
	"fmt"
	"github.com/jeanschmitt/tzgen/examples/contracts"
	"log"
	"os"
)

//go:generate go run ../main.go --address KT1K3ZqbYq1bCwpSPNX9xBgQd8CaYxRVXd4P --name Hello --pkg contracts -o ./contracts/Hello.go

var examples = map[string]func(){
	"storage": Storage,
}

func Storage() {
	hello := helloContract(rpcClient())

	value, err := hello.Storage(context.Background())
	handleErr(err)

	fmt.Printf("Current storage: %v\n", value)
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
	tzClient, err := rpc.NewClient("https://ghostnet.smartpy.io", nil)
	handleErr(err)

	err = tzClient.Init(context.Background())
	handleErr(err)

	tzClient.Listen()
	handleErr(err)

	return tzClient
}

func helloContract(tzClient *rpc.Client) *contracts.Hello {
	return contracts.NewHello(tezos.MustParseAddress("KT1K3ZqbYq1bCwpSPNX9xBgQd8CaYxRVXd4P"), tzClient)
}

func handleErr(err error) {
	if err != nil {
		log.Fatalf("An error occured: %v\n", err)
	}
}
