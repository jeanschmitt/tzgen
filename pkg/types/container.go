package types

// Struct is an aggregation of named fields.
// It corresponds to non-top-level `pair` prims, in Michelson.
// It can represent either a parameter's type of entrypoint, or a record in a
// storage type.
type Struct struct {
	Name   string
	Fields []Param
}

func (s *Struct) Hash() []byte {
	// Name should not be part of the hash payload
	hashes := [][]byte{[]byte(TypeStruct)}
	for _, f := range s.Fields {
		hashes = append(hashes, f.Hash())
	}
	return TypeHash(hashes...)
}

func (Struct) TypeName() string {
	return TypeStruct
}

// Union is a type that can be either its Left type, or its Right one.
// It corresponds to `or` types in Michelson.
type Union struct {
	Name  string
	Left  Type
	Right Type
}

func (u *Union) Hash() []byte {
	return TypeHash([]byte(TypeUnion), u.Left.Hash(), u.Right.Hash())
}

func (Union) TypeName() string {
	return TypeUnion
}

type List struct {
	Type Type
}

func (l *List) Hash() []byte {
	return TypeHash([]byte(TypeList), l.Type.Hash())
}

func (List) TypeName() string {
	return TypeList
}

type Set struct {
	Type Type
}

func (s *Set) Hash() []byte {
	return TypeHash([]byte(TypeSet), s.Type.Hash())
}

func (Set) TypeName() string {
	return TypeSet
}

type Map struct {
	Key   Type
	Value Type
}

func (m *Map) Hash() []byte {
	return TypeHash([]byte(TypeMap), m.Key.Hash(), m.Value.Hash())
}

func (Map) TypeName() string {
	return TypeMap
}

type BigMap struct {
	Key   Type
	Value Type
}

func (b *BigMap) Hash() []byte {
	return TypeHash([]byte(TypeBigmap), b.Key.Hash(), b.Value.Hash())
}

func (BigMap) TypeName() string {
	return TypeBigmap
}

type Option struct {
	Type Type
}

func (o *Option) Hash() []byte {
	return TypeHash([]byte(TypeOption), o.Type.Hash())
}

func (Option) TypeName() string {
	return TypeOption
}
