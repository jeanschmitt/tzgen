package codegen

import "strings"

type normalizedName struct {
	Original   string
	Normalized string
}

type goFormat struct{}

var GoFormat = goFormat{}

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

func normalize(original string, normalizeFunc func(string) string) normalizedName {
	return normalizedName{
		Normalized: normalizeFunc(original),
		Original:   original,
	}
}

func ToPascalCase(input string) string {
	parts := strings.Split(input, "_")
	for i, s := range parts {
		if len(s) > 0 {
			parts[i] = strings.ToUpper(s[:1]) + s[1:]
		}
	}
	return strings.Join(parts, "")
}

// TODO
func ToCamelCase(input string) string {
	return input
}
