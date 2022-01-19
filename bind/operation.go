package bind

import (
	"context"

	"github.com/completium/go-tezos/v4/rpc"
)

type Operation struct {
	Destination string
	Entrypoint  string
	Value       string
	*PreApplyRes
}

func NewOperation(destination, entrypoint, value string) *Operation {
	return &Operation{
		Destination: destination,
		Entrypoint:  entrypoint,
		Value:       value,
	}
}

type PreApplyRes struct {
}

type Injection struct {
	*Operation
	Hash string
}

func (o *Operation) PreApplyAndInject(ctx context.Context, tz *rpc.Client) (*Injection, error) {
	o, err := o.PreApply(ctx, tz)
	if err != nil {
		return nil, err
	}

	return o.Inject(ctx, tz)
}

func (o *Operation) PreApply(ctx context.Context, tz *rpc.Client) (*Operation, error) {
	o.PreApplyRes = &PreApplyRes{}
	return o, nil
}

func (o *Operation) Inject(ctx context.Context, tz *rpc.Client) (*Injection, error) {
	return &Injection{
		Operation: o,
		Hash:      "4242",
	}, nil
}
