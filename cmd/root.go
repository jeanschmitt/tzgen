package cmd

import (
	"bytes"
	"github.com/iancoleman/strcase"
	"github.com/jeanschmitt/tzgen/internal/generate"
	"github.com/jeanschmitt/tzgen/internal/parse"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tzgen",
	Short:   "Generate Go bindings to a Tezos smart contract",
	Version: Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		src, err := getSrc()
		if err != nil {
			return errors.Wrap(err, "failed to get contract script")
		}

		generated, err := generateBindings(src)
		if err != nil {
			return errors.Wrap(err, "failed to generate bindings")
		}

		err = writeResult(generated)
		if err != nil {
			return errors.Wrap(err, "failed to write generated code to file")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

const (
	defaultEndpoint = "https://ghostnet.smartpy.io"
)

var (
	endpointFlag  string
	addressFlag   string
	srcFlag       string
	nameFlag      string
	pkgFlag       string
	outFlag       string
	fixupFileFlag string
)

func init() {
	rootCmd.Flags().StringVarP(&endpointFlag, "endpoint", "e", defaultEndpoint, "rpc endpoint to use")
	rootCmd.Flags().StringVarP(&addressFlag, "address", "a", "", "address of the contract. Required if --src is not set")
	rootCmd.Flags().StringVar(&srcFlag, "src", "", "json file containing the contract's script")
	_ = rootCmd.MarkFlagFilename("src")

	rootCmd.Flags().StringVar(&nameFlag, "name", "", "name of the contract")
	rootCmd.Flags().StringVar(&pkgFlag, "pkg", "", "go package of the output code")
	rootCmd.Flags().StringVarP(&outFlag, "out", "o", "", "output file. Prints to Stdout if not set")
	rootCmd.Flags().StringVarP(&fixupFileFlag, "fixup", "f", "", "yaml fixup file to use")
	_ = rootCmd.MarkFlagRequired("name")
	_ = rootCmd.MarkFlagRequired("pkg")
	_ = rootCmd.MarkFlagFilename("out", "go")
	_ = rootCmd.MarkFlagFilename("fixup")
}

func generateBindings(script []byte) ([]byte, error) {
	var err error
	data := generate.Data{
		Address: addressFlag,
		Package: pkgFlag,
	}
	data.Contract, data.Structs, err = parse.Parse(script, nameFlag)
	if err != nil {
		return nil, err
	}

	if fixupFileFlag != "" {
		fixupFile, err := os.ReadFile(fixupFileFlag)
		if err != nil {
			return nil, err
		}

		var fixupCfg parse.FixupConfig
		err = yaml.NewDecoder(bytes.NewReader(fixupFile)).Decode(&fixupCfg)
		if err != nil {
			return nil, err
		}

		data.Structs = parse.Fixup(fixupCfg, data.Structs, strcase.ToCamel)
	}

	return generate.Render(&data)
}

func getSrc() ([]byte, error) {
	if srcFlag != "" {
		return os.ReadFile(srcFlag)
	}

	// Get source from RPC
	// At this point, addressFlag is required
	if addressFlag == "" {
		return nil, errors.New("--address is required when getting script from rpc")
	}

	u, err := url.JoinPath(endpointFlag, "chains/main/blocks/head/context/contracts", addressFlag, "script")
	if err != nil {
		return nil, err
	}
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func writeResult(out []byte) error {
	if outFlag == "" {
		_, err := os.Stdout.Write(out)
		if err != nil {
			return err
		}
		return nil
	}

	return os.WriteFile(outFlag, out, 0o644)
}
