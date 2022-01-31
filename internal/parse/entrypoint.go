package parse

import (
	"blockwatch.cc/tzgo/micheline"
	"fmt"
	types "github.com/jeanschmitt/tzgen/pkg/types"
	"github.com/pkg/errors"
	"sort"
)

func (p *Parser) entrypoints() error {
	tzgoEntrypoints, err := p.script.Entrypoints(true)
	if err != nil {
		return err
	}

	p.contract.Entrypoints = make([]Entrypoint, 0, len(tzgoEntrypoints))
	for name, e := range tzgoEntrypoints {
		parsed, err := p.entrypoint(name, &e)
		if err != nil {
			return errors.Wrapf(err, "failed to parse `%s` entrypoint", name)
		}
		if parsed != nil {
			p.contract.Entrypoints = append(p.contract.Entrypoints, *parsed)
		}
	}

	sort.SliceStable(p.contract.Entrypoints, func(i, j int) bool {
		return p.contract.Entrypoints[i].id < p.contract.Entrypoints[j].id
	})

	return nil
}

func (p *Parser) entrypoint(name string, e *micheline.Entrypoint) (*Entrypoint, error) {
	entrypoint := Entrypoint{
		Name: name,
		id:   e.Id,
	}

	for _, arg := range e.Typedef {
		typ, err := p.parseType(&arg)
		if err != nil {
			return nil, err
		}
		switch typ.(type) {
		case types.Contract:
			fmt.Printf("Skipping `%s` getter\n", name)
			return nil, nil
		}
		entrypoint.Args = append(entrypoint.Args, types.Param{Name: arg.Name, Type: typ})
	}

	return &entrypoint, nil
}
