package main

import (
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
)

// Generate bindings used in examples

//go:generate go run github.com/jeanschmitt/tzgen --name Hello -o ./contracts/Hello.go -a KT1FsGDhS7PUWn1Yff8r5nXE4MArZwD2XhEi --pkg contracts -e https://ghostnet.smartpy.io
//go:generate go run github.com/jeanschmitt/tzgen --name Payable -o ./contracts/Payable.go -a KT1FPA7vN4cBk24df7VxUu9DstRcc7am3qnf --pkg contracts -e https://ghostnet.smartpy.io
//go:generate go run github.com/jeanschmitt/tzgen --name FA2NFT -o ./contracts/FA2NFT.go -a KT1UvFD9hwgbLKbmTzUPmVRVsHG29ciQrcwE --pkg contracts -e https://ghostnet.smartpy.io

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: go run github.com/jeanschmitt/tzgen/examples [examples]")
		fmt.Println("Available examples:")
		for _, ex := range exampleList() {
			fmt.Printf("\t%s\n", ex)
		}

		os.Exit(0)
	}

	for _, ex := range flag.Args() {
		fmt.Printf("Example %q\n", ex)
		runExample(ex)
		fmt.Println()
	}
}

type Examples struct{}

var examples Examples

func exampleList() []string {
	examplesType := reflect.ValueOf(examples).Type()
	var methods []string

	for i := 0; i < examplesType.NumMethod(); i++ {
		methods = append(methods, examplesType.Method(i).Name)
	}

	return methods
}

func runExample(name string) {
	examplesVal := reflect.ValueOf(examples)

	methodVal := examplesVal.MethodByName(name)
	if methodVal.IsZero() {
		log.Fatalf("Example %q doesn't exist\n", name)
	}

	methodVal.Call(nil)
}

const (
	Mainnet  = "https://mainnet.smartpy.io"
	Ghostnet = "https://ghostnet.smartpy.io"
)

var ctx = context.Background()

func NewClient(rpcUrl string) *rpc.Client {
	c, err := rpc.NewClient(rpcUrl, nil)
	fatalOnErr(err)

	err = c.Init(ctx)
	fatalOnErr(err)

	c.Listen()

	return c
}

var AlicePrivateKey = tezos.MustParsePrivateKey("edsk3QoqBuvdamxouPhin7swCvkQNgq4jP5KZPbwWNnwdZpSpJiEbq")

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal("Fatal error: ", err)
	}
}

func must[T any](t T, err error) T {
	fatalOnErr(err)
	return t
}
