# Type mapping

Micheline <=> Go type mapping

## Core data types

| micheline | go        |
|-----------|-----------|
| string    | string    |
| nat       | *big.Int  |
| int       | *big.Int  |
| bytes     | []byte    |
| bool      | bool      |
| unit      | (nothing) |
| never     | (nothing) |

## Container types

| micheline   | go             |
|-------------|----------------|
| list T      | []T            |
| pair L R    | struct         |
| option T    | bind.Option[T] |
| or L R      | bind.Or[L, R]  |
| set T       | []T            |
| map K T     | TODO           |
| big_map K T | TODO           |

## Domain types

| micheline    | readable         | optimized | go        |
|--------------|------------------|-----------|-----------|
| mutez        | nat              | nat       | *big.Int  |
| timestamp    | string (RFC3339) | nat       | time.Time |
| contract     | string (Base58)  | bytes     | string    |
| address      | string (Base58)  | bytes     | string    |
| key          | string (Base58)  | bytes     | string    |
| signature    | string (Base58)  | bytes     | string    |
| bls12_381_g1 | bytes            | bytes     | []byte    |
| bls12_381_g2 | bytes            | bytes     | []byte    |
| bls12_381_fr | bytes            | bytes     | []byte    |
