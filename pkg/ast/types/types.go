//go:generate go run gen.go -out builtin.go Nat Int String Bool Bytes Unit Timestamp Address Mutez Key KeyHash Signature ChainID Operation Contract

package types

import (
	"crypto/md5"
)

type Type interface {
	TypeName() string
	Hash() []byte
}

// Builtin types
const (
	TypeNat       = "nat"
	TypeInt       = "int"
	TypeString    = "string"
	TypeBool      = "bool"
	TypeBytes     = "bytes"
	TypeUnit      = "unit"
	TypeTimestamp = "timestamp"
	TypeAddress   = "address"
	TypeMutez     = "mutez"
	TypeKey       = "key"
	TypeKeyHash   = "key_hash"
	TypeSignature = "signature"
	TypeChainID   = "chain_id"
	TypeOperation = "operation"
	TypeContract  = "contract"
)

// Container types
const (
	TypeStruct = "struct"
	TypeUnion  = "union"
	TypeList   = "list"
	TypeSet    = "set"
	TypeMap    = "map"
	TypeBigmap = "big_map"
	TypeOption = "option"
)

const HashSize = md5.Size

func TypeHash(elems ...[]byte) []byte {
	hash := md5.New()
	for _, e := range elems {
		hash.Write(e)
	}
	return hash.Sum(nil)
}

// Param is a type associated with a name.
// It can be an argument of an entrypoint or a field of a struct.
type Param struct {
	Name         string
	OriginalType string // TODO: keep?
	Type         Type
}

func (p Param) Hash() []byte {
	hash := md5.New()
	hash.Write([]byte(p.Name))
	hash.Write(p.Type.Hash())
	return hash.Sum(nil)
}
