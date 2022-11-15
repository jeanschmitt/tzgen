package types

// Struct is an aggregation of named fields.
// It corresponds to non-top-level `pair` prims, in Michelson.
// It can represent either a parameter's type of entrypoint, or a record in a
// storage type.
type Struct struct {
	Name   string
	Fields []Param
	// If true, the expected prim matching to this struct has a flat structure,
	// instead of a tree of pairs.
	Flat bool
}

func (s *Struct) Hash() []byte {
	hashes := [][]byte{[]byte(TypeStruct + s.Name)}
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

type Bigmap struct {
	Key   Type
	Value Type
}

func (b *Bigmap) Hash() []byte {
	return TypeHash([]byte(TypeBigmap), b.Key.Hash(), b.Value.Hash())
}

func (Bigmap) TypeName() string {
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

type Lambda struct {
	Param  Type
	Return Type
}

func (l *Lambda) Hash() []byte {
	return TypeHash([]byte(TypeLambda), l.Param.Hash(), l.Return.Hash())
}

func (Lambda) TypeName() string {
	return TypeLambda
}
