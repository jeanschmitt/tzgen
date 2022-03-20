# tzgen

Go binding to Tezos smart contracts, using code generation.

## Installation

```bash
go install github.com/jeanschmitt/tzgen@latest
```

## Usage

### From a deployed contract

```bash
tzgen --name Simple --pkg contracts --address KT1CiYNu9iJknnL31TXBWHCqRdFRh7jPWdzg -o ./contracts/Simple.go
```

The endpoint and chain id are `https://hangzhounet.smartpy.io` and `main` by default, but can be overridden with `--node` and `--chain`.

### From a micheline file

```bash
tzgen --name Simple --pkg contracts -i ./Simple.json -o ./contracts/Simple.go
```

## Note

This tool is still under development, so the generated bindings don't cover every smart contract yet.
