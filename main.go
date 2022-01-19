package main

import (
	"fmt"
	"os"
	"path"

	flag "github.com/spf13/pflag"

	"github.com/jeanschmitt/tzgen/bind/codegen"
)

func main() {
	inFile, outFile, pkgName, contract := getOptions()

	micheline, err := os.ReadFile(inFile)
	if err != nil {
		panic(err)
	}

	out, err := codegen.Generate(micheline, pkgName, contract)
	if err != nil {
		panic(err)
	}

	outF, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}

	n, err := outF.Write(out)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d bytes written\n", n)
}

var forbiddenInputExt = []string{
	".go",
}

func getOptions() (input, goOutput, pkgName, contract string) {
	goOut := flag.String("out", "", "path for generated file")
	packageName := flag.String("pkg", "", "name of the go package to use")
	contractName := flag.String("name", "", "name of the contract")
	_ = flag.BoolP("verbose", "v", false, "")
	help := flag.BoolP("help", "h", false, "print help")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(1)
	}

	if flag.NArg() != 1 {
		if flag.NArg() == 0 {
			exitBadArgs("Input file is missing")
		} else {
			exitBadArgs("Unexpected positional parameter")
		}
	}
	if *goOut == "" {
		exitBadArgs("Output file is missing")
	}
	if *packageName == "" {
		exitBadArgs("Package name is missing")
	}
	if *contractName == "" {
		exitBadArgs("Contract name is missing")
	}

	inputFile := flag.Arg(0)
	inputExt := path.Ext(inputFile)

	if stringInSlice(inputExt, forbiddenInputExt) {
		fmt.Printf("Forbidden input extension: %s\n", inputExt)
		os.Exit(1)
	}

	return inputFile, *goOut, *packageName, *contractName
}

func printHelp() {
	fmt.Println(`Usage:
	tzgen --out <out file> --pkg <package> --name <contract name> <input file>
	`)
	flag.PrintDefaults()
}

func exitBadArgs(reason string) {
	fmt.Println(reason)
	printHelp()
	os.Exit(1)
}

func stringInSlice(str string, slice []string) bool {
	for _, i := range slice {
		if str == i {
			return true
		}
	}
	return false
}
