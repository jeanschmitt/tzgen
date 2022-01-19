package codegen

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplate(t *testing.T) {
	union := UnionType{
		Name:  normalizedName{Normalized: "Unknown"},
		RType: &StringType{},
	}
	union.LType = &union

	structT := StructType{
		Name: normalizedName{Normalized: "AStruct"},
		Fields: []*Field{
			{
				Name: normalizedName{Normalized: "TheUnion"},
				Type: &union,
			},
			{
				Name: normalizedName{Normalized: "nickname"},
				Type: &StringType{},
			},
		},
	}

	data := Data{
		Package: "template",
		Contract: Contract{
			Name:      "Quartz",
			Micheline: "micheline",
			Entrypoints: []*Entrypoint{
				{
					Name: normalizedName{Normalized: "Mint", Original: "mint"},
					Args: []*Arg{
						{
							Name: normalizedName{Normalized: "to"},
							Type: &StringType{},
						},
						{
							Name: normalizedName{Normalized: "amount"},
							Type: &IntType{},
						},
						{
							Name: normalizedName{Normalized: "tokenID"},
							Type: &IntType{},
						},
					},
				},
				{
					Name: normalizedName{Normalized: "UnionNous", Original: "union_nous"},
					Args: []*Arg{
						{
							Name: normalizedName{Normalized: "un"},
							Type: &union,
						},
						{
							Name: normalizedName{Normalized: "structs"},
							Type: &ListType{T: &structT},
						},
					},
				},
			},
		},
		Unions: []*UnionType{
			&union,
		},
		Structs: []*StructType{
			&structT,
		},
	}

	out, err := data.Render(true)
	require.NoError(t, err)

	f, _ := os.Create("out.go")
	f.Write(out)
}
