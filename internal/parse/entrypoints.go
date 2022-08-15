package parse

import (
	"blockwatch.cc/tzgo/micheline"
	"github.com/jeanschmitt/tzgen/pkg/ast"
	"github.com/jeanschmitt/tzgen/pkg/ast/types"
	"github.com/pkg/errors"
	"sort"
	"strconv"
)

func (p *parser) parseEntrypoints() error {
	entrypointMap, err := p.script.Entrypoints(true)
	if err != nil {
		return errors.Wrap(err, "failed to get entrypoints")
	}

	entrypoints := mapValues(entrypointMap)

	// Sort entrypoints by id
	sort.SliceStable(entrypoints, func(i, j int) bool { return entrypoints[i].Id < entrypoints[j].Id })

	for _, entrypoint := range entrypoints {
		e, err := p.parseEntrypoint(&entrypoint)
		if err != nil {
			return errors.Wrap(err, "failed to parse entrypoint")
		}
		if getter, isGetter := e.(*ast.Getter); isGetter {
			p.contract.Getters = append(p.contract.Getters, getter)
		} else {
			p.contract.Entrypoints = append(p.contract.Entrypoints, e.(*ast.Entrypoint))
		}
	}

	return nil
}

func (p *parser) parseEntrypoint(entrypoint *micheline.Entrypoint) (any, error) {
	e := ast.Entrypoint{
		Name: entrypoint.Name,
		Raw:  entrypoint,
	}

	nArgs := len(entrypoint.Typedef)
	for i, arg := range entrypoint.Typedef {
		if arg.Type == "unit" && i == 0 {
			// continue because it can still be a getter
			continue
		}
		if arg.Type == "contract" && i == nArgs-1 {
			// arg.Args contains the return type of the getter
			returnType, err := p.parseType(&arg.Args[0])
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse return type")
			}
			return &ast.Getter{Entrypoint: e, ReturnType: returnType}, nil
		}

		typ, err := p.parseType(&arg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse type")
		}

		e.Params = append(e.Params, entrypointParam(&arg, typ, i))
	}

	return &e, nil
}

func entrypointParam(arg *micheline.Typedef, typ types.Type, i int) *types.Param {
	argName := arg.Name
	if argName == "" || startsWithInt(argName) {
		argName = arg.Type + strconv.Itoa(i)
	}
	originalType := arg.Type
	if arg.Optional {
		originalType = "option<" + originalType + ">"
	}

	return &types.Param{Name: argName, Type: typ, OriginalType: originalType}
}

func startsWithInt(s string) bool {
	if s == "" {
		return false
	}
	return s[0] >= '0' && s[0] <= '9'
}

func mapValues[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
