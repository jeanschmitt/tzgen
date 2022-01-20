package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/jeanschmitt/tzgen/internal/codegen"
)

// CLI args
var (
	outPath      string
	pkgName      string
	contractName string
	verbose      bool
)

func main() {
	parseFlags()
	configLog(verbose)

	micheline, err := getInput()
	if err != nil {
		logFatal(err)
	}

	out, err := codegen.Generate(micheline, pkgName, contractName)
	if err != nil {
		logFatal(err)
	}

	if err := writeOutput(out); err != nil {
		logFatal(err)
	}
}

func parseFlags() {
	flag.StringVar(&outPath, "out", "", "path for the generated file")
	flag.StringVar(&pkgName, "pkg", "", "name of the go package to use")
	flag.StringVar(&contractName, "name", "", "name of the contract")
	flag.BoolVarP(&verbose, "verbose", "v", false, "")
	help := flag.BoolP("help", "h", false, "print help")
	version := flag.Bool("version", false, "print version")

	flag.Parse()

	if *help || flag.NArg()+flag.NFlag() == 0 {
		printUsage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("tzgen version %s\n", Version)
		os.Exit(0)
	}

	if pkgName == "" {
		exitBadArgs("Package name is missing")
	}
	if contractName == "" {
		exitBadArgs("Contract name is missing")
	}
}

func getInput() (input []byte, err error) {
	fi, _ := os.Stdin.Stat()

	if fi.Mode()&os.ModeCharDevice == 0 {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read from standard input")
		}
		return input, nil
	}

	// If we are not using piped input, the micheline code is read from a file provided
	// as cli argument.

	if flag.NArg() != 1 {
		if flag.NArg() == 0 {
			exitBadArgs("Input file is missing")
		} else {
			exitBadArgs("Unexpected positional parameter")
		}
	}

	inputFile := flag.Arg(0)
	inputExt := path.Ext(inputFile)

	if stringInSlice(inputExt, forbiddenInputExt) {
		return nil, errors.Errorf("Forbidden input extension: %s\n", inputExt)
	}

	micheline, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read input file")
	}

	return micheline, nil
}

func writeOutput(output []byte) error {
	if outPath == "" {
		_, err := os.Stdout.Write(output)
		return err
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	n, err := outFile.Write(output)
	if err != nil {
		return err
	}

	log.Infof("%d bytes written\n", n)

	return nil
}

func configLog(verbose bool) {
	log.SetFormatter(&log.TextFormatter{})

	level := log.InfoLevel
	if verbose {
		level = log.TraceLevel
	}
	log.SetLevel(level)
}

func logFatal(err error) {
	if err == nil {
		return
	}
	log.Fatalf("Fatal error: %v", err)
}

var forbiddenInputExt = []string{
	".go",
}

func printUsage() {
	fmt.Println(`Usage:
	tzgen --out <out file> --pkg <package> --name <contract name> <input file>
	`)
	flag.PrintDefaults()
}

func exitBadArgs(reason string) {
	fmt.Println(reason)
	printUsage()
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
