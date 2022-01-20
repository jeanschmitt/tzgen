package codegen

import "strings"

type Formatter interface {
	Method(original string) normalizedName
	Struct(original string) normalizedName
	Argument(original string) normalizedName
	Field(original string) normalizedName
}

var (
	GoFormat = goFormat{}
)

type normalizedName struct {
	Original   string
	Normalized string
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
