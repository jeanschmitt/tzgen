package codegen

type Type interface {
	IsRefType() bool
	GoType() string
	GoPrim(arg string) string
}

type baseType struct{}

func (t baseType) IsRefType() bool {
	return false
}

type IntType struct{ baseType }

func (t *IntType) GoType() string {
	return "*big.Int"
}

func (t *IntType) GoPrim(arg string) string {
	return "Int(" + arg + ")"
}

type StringType struct{ baseType }

func (t *StringType) GoType() string {
	return "string"
}

func (t *StringType) GoPrim(arg string) string {
	return "String(" + arg + ")"
}

type BytesType struct{ baseType }

func (t *BytesType) GoType() string {
	return "[]byte"
}

func (t *BytesType) GoPrim(arg string) string {
	return "Bytes(" + arg + ")"
}

type ListType struct {
	baseType
	T Type
}

func (t *ListType) GoType() string {
	return "[]" + t.T.GoType()
}

func (t *ListType) GoPrim(arg string) string {
	return "List(" + arg + ")"
}

type MapType struct {
	baseType
	Key   Type
	Value Type
}

func (t *MapType) GoType() string {
	return "map[" + t.Key.GoType() + "]" + t.Value.GoType()
}

type OptionType struct {
	baseType
	Name normalizedName
	Type Type
}

func (t *OptionType) GoPrim(arg string) string {
	return "Primer(" + arg + ")"
}

func (t *OptionType) GoType() string {
	return "*" + t.Name.Normalized
}

func (t *OptionType) IsRefType() bool {
	return true
}

type StructType struct {
	baseType
	Name   normalizedName
	Fields []*Field
}

func (t *StructType) GoPrim(arg string) string {
	return "Primer(" + arg + ")"
}

func (t *StructType) GoType() string {
	return "*" + t.Name.Normalized
}

func (t *StructType) IsRefType() bool {
	return true
}

type UnionType struct {
	baseType
	Name  normalizedName
	LType Type
	RType Type
}

func (t *UnionType) GoPrim(arg string) string {
	return "Primer(" + arg + ")"
}

func (t *UnionType) GoType() string {
	return "*" + t.Name.Normalized
}

func (t *UnionType) IsRefType() bool {
	return true
}
