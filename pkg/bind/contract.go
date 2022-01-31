package bind

import (
	"context"
	"net/http"

	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/pkg/errors"
)

type Contract struct {
	Address tezos.Address
	Rpc     *Client
}

func NewContract(address string, c *Client) (*Contract, error) {
	tzAddress, err := tezos.ParseAddress(address)
	if err != nil {
		return nil, errors.Wrap(err, "invalid address")
	}
	return &Contract{Address: tzAddress, Rpc: c}, nil
}

// IsDeployed returns wether the contract is deployed at its configured address or not.
func (c *Contract) IsDeployed(ctx context.Context) (bool, error) {
	_, err := c.Rpc.GetContractStorage(ctx, c.Address, rpc.Head)
	if err != nil {
		var httpError rpc.HTTPError
		if errors.As(err, &httpError) && httpError.StatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
