package types

type Nat struct{}

func (Nat) Hash() []byte {
	return TypeHash([]byte(TypeNat))
}

func (Nat) TypeName() string {
	return TypeNat
}

type Int struct{}

func (Int) Hash() []byte {
	return TypeHash([]byte(TypeInt))
}

func (Int) TypeName() string {
	return TypeInt
}

type String struct{}

func (String) Hash() []byte {
	return TypeHash([]byte(TypeString))
}

func (String) TypeName() string {
	return TypeString
}

type Bool struct{}

func (Bool) Hash() []byte {
	return TypeHash([]byte(TypeBool))
}

func (Bool) TypeName() string {
	return TypeBool
}

type Bytes struct{}

func (Bytes) Hash() []byte {
	return TypeHash([]byte(TypeBytes))
}

func (Bytes) TypeName() string {
	return TypeBytes
}

type Unit struct{}

func (Unit) Hash() []byte {
	return TypeHash([]byte(TypeUnit))
}

func (Unit) TypeName() string {
	return TypeUnit
}

type Timestamp struct{}

func (Timestamp) Hash() []byte {
	return TypeHash([]byte(TypeTimestamp))
}

func (Timestamp) TypeName() string {
	return TypeTimestamp
}

type Address struct{}

func (Address) Hash() []byte {
	return TypeHash([]byte(TypeAddress))
}

func (Address) TypeName() string {
	return TypeAddress
}

type Mutez struct{}

func (Mutez) Hash() []byte {
	return TypeHash([]byte(TypeMutez))
}

func (Mutez) TypeName() string {
	return TypeMutez
}

type Key struct{}

func (Key) Hash() []byte {
	return TypeHash([]byte(TypeKey))
}

func (Key) TypeName() string {
	return TypeKey
}

type KeyHash struct{}

func (KeyHash) Hash() []byte {
	return TypeHash([]byte(TypeKeyHash))
}

func (KeyHash) TypeName() string {
	return TypeKeyHash
}

type Signature struct{}

func (Signature) Hash() []byte {
	return TypeHash([]byte(TypeSignature))
}

func (Signature) TypeName() string {
	return TypeSignature
}

type ChainID struct{}

func (ChainID) Hash() []byte {
	return TypeHash([]byte(TypeChainID))
}

func (ChainID) TypeName() string {
	return TypeChainID
}

type Operation struct{}

func (Operation) Hash() []byte {
	return TypeHash([]byte(TypeOperation))
}

func (Operation) TypeName() string {
	return TypeOperation
}

type Contract struct{}

func (Contract) Hash() []byte {
	return TypeHash([]byte(TypeContract))
}

func (Contract) TypeName() string {
	return TypeContract
}