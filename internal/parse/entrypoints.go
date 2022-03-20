package parse

import (
	"blockwatch.cc/tzgo/micheline"
	"github.com/jeanschmitt/tzgen/pkg/ast"
	"github.com/jeanschmitt/tzgen/pkg/ast/types"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"sort"
	"strconv"
)

func (p *parser) parseEntrypoints() error {
	entrypointMap, err := p.script.Entrypoints(true)
	if err != nil {
		return errors.Wrap(err, "failed to get entrypoints")
	}

	entrypoints := maps.Values(entrypointMap)

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

func (p *parser) parseEntrypoint(entrypoint *micheline.Entrypoint) (interface{}, error) {
	e := ast.Entrypoint{
		Name: entrypoint.Name,
		Raw:  entrypoint,
	}

	nArgs := len(entrypoint.Typedef)
	_ = nArgs
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

		argName := arg.Name
		if argName == "" || startsWithInt(argName) {
			argName = arg.Type + strconv.Itoa(i)
		}
		originalType := arg.Type
		if arg.Optional {
			originalType = "option<" + originalType + ">"
		}
		e.Params = append(e.Params, &types.Param{Name: argName, Type: typ, OriginalType: originalType})
	}

	return &e, nil
}

func startsWithInt(s string) bool {
	if s == "" {
		return false
	}
	return s[0] >= '0' && s[0] <= '9'
}
