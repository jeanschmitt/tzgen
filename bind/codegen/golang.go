package codegen

type goFormat struct{}



func (g goFormat) Method(original string) normalizedName {
	return normalize(original, ToPascalCase)
}

func (g goFormat) Struct(original string) normalizedName {
	return normalize(original, ToPascalCase)
}

func (g goFormat) Argument(original string) normalizedName {
	return normalize(original, ToCamelCase)
}

func (g goFormat) Field(original string) normalizedName {
	return normalize(original, ToCamelCase)
}