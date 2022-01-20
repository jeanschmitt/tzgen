# tzgen

Go binding to Tezos smart contracts, using code generation.

## Build

```bash
make build
```

## Usage

```bash
tzgen --out <out file> --pkg <package> --name <contract name> <input file>
```

It can also read from standard input:

```bash
curl https://rpc.tzkt.io/mainnet/chains/main/blocks/head/context/contracts/KT1FvqJwEDWb1Gwc55Jd1jjTHRVWbYKUUpyq/script | tzgen --out <out file> --pkg <package> --name <contract name>
```
