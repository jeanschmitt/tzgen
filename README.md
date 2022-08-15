# tzgen

Go binding to Tezos smart contracts, using code generation.

## Installation

```bash
go install github.com/jeanschmitt/tzgen@latest
```

## Usage

### From a deployed contract

```bash
tzgen --name Hello --pkg contracts --address KT1K3ZqbYq1bCwpSPNX9xBgQd8CaYxRVXd4P -o ./contracts/Hello.go
```

The endpoint is `https://ghostnet.smartpy.io` by default, but can be overridden with `--endpoint`.

### From a micheline file

```bash
tzgen --name Hello --pkg contracts --src ./Hello.json -o ./contracts/Hello.go
```
