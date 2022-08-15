package main

import (
	_ "embed"
	"fmt"
	"github.com/jeanschmitt/tzgen/internal/generate"
	"github.com/jeanschmitt/tzgen/internal/parse"
	flag "github.com/spf13/pflag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// input
var (
	nodeURL = flag.String("node", "https://ghostnet.smartpy.io", "endpoint to use when using --address")
	chain   = flag.String("chain", "main", "chain id to use when using --address")
	address = flag.String("address", "", "address of the deployed contract")
	inFile  = flag.StringP("infile", "i", "", "path for the micheline input file")
)

// output
var (
	contractName = flag.String("name", "", "name of the contract (and of the generated struct type)")
	packageName  = flag.String("pkg", "", "package of the generated go file")
	outFile      = flag.StringP("outfile", "o", "", "path to the output file. Prints to standard output if not provided")
)

// others
var (
	helpFlag    = flag.BoolP("help", "h", false, "print help")
	versionFlag = flag.BoolP("version", "V", false, "show version")
)

//go:embed VERSION
var Version string

func main() {
	if len(os.Args) <= 1 {
		usage()
		os.Exit(0)
	}
	flag.Parse()
	if *helpFlag {
		usage()
		os.Exit(0)
	}
	if *versionFlag {
		fmt.Println("tzgen version", Version)
		os.Exit(0)
	}
	validateArgs()

	script := getInput()

	contract, structs, unions, err := parse.Parse(script, *contractName)
	handleErr(err)

	data := &generate.Data{
		Contract: contract,
		Structs:  structs,
		Unions:   unions,
		Address:  *address,
		Package:  *packageName,
	}

	out, err := generate.Render(data)
	handleErr(err)

	processOutput(out)
}

func validateArgs() {
	if *inFile != "" && *address != "" {
		badArgs("--address and --infile are mutually exclusive")
	}
	if *inFile == "" && *address == "" {
		badArgs("no input provided: please provide --address or --infile")
	}
	if *contractName == "" {
		// Note: we may implement name resolution with TZIP-16 in the future
		badArgs("--name is required")
	}
	if *packageName == "" {
		// Note: it could be guessed from directory
		badArgs("--pkg is required")
	}
}

func getInput() []byte {
	if *inFile != "" {
		data, err := os.ReadFile(*inFile)
		handleErr(err)

		return data
	}
	if *address != "" {
		res, err := http.Get(scriptURL())
		handleErr(err)
		defer func() { _ = res.Body.Close() }()

		data, err := io.ReadAll(res.Body)
		handleErr(err)
		return data
	}

	panic("args were not validated")
}

func processOutput(out []byte) {
	if *outFile == "" {
		_, err := os.Stdout.Write(out)
		handleErr(err)
		return
	}

	err := os.WriteFile(*outFile, out, 0644)
	handleErr(err)
}

func scriptURL() string {
	u, err := url.Parse(*nodeURL)
	handleErr(err)

	path := fmt.Sprintf("chains/%s/blocks/head/context/contracts/%s/script", *chain, *address)
	u, err = u.Parse(path)
	handleErr(err)

	return u.String()
}

func usage() {
	fmt.Println(`Usage:
    tzgen [--help] [--version] [--outfile] [--name] [--pkg]
          [--infile | --address] [--node] [--chain]

Options:`)
	flag.PrintDefaults()
}

func badArgs(msg string) {
	fmt.Println("Invalid arguments:", msg)
	usage()
	os.Exit(0)
}

func handleErr(err error) {
	if err != nil {
		log.Fatal("Fatal error:", err)
	}
}
