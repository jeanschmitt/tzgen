package parse

import "github.com/jeanschmitt/tzgen/pkg/ast/types"

type FixupConfig map[string]FixupStruct

type FixupStruct struct {
	Name   string
	Equals string
	Fields map[string]string
}

func Fixup(cfg FixupConfig, structs []*types.Struct, processNameFunc func(string) string) []*types.Struct {
	structsByOldName := make(map[string]*types.Struct)
	for _, s := range structs {
		structsByOldName[processNameFunc(s.Name)] = s
	}

	for name, v := range cfg {
		currentStruct := structsByOldName[name]

		if v.Name != "" {
			currentStruct.Name = v.Name
		}

		for oldF, newF := range v.Fields {
			for i, f := range currentStruct.Fields {
				if processNameFunc(f.Name) == oldF {
					currentStruct.Fields[i].Name = newF
				}
			}
		}
	}

	for name, v := range cfg {
		old := structsByOldName[name]
		if v.Equals != "" {
			*old = *structsByOldName[v.Equals]
			structs = deleteStruct(structs, old)
		}
	}

	return structs
}

func deleteStruct(structs []*types.Struct, str *types.Struct) []*types.Struct {
	for i, s := range structs {
		if s == str {
			newSlice := make([]*types.Struct, 0, len(structs)-1)
			newSlice = append(newSlice, structs[:i]...)
			return append(newSlice, structs[i+1:]...)
		}
	}
	return structs
}
