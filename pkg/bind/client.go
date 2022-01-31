package bind

import (
	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Client struct {
	*rpc.Client
	Constants rpc.Constants
	Protocol  tezos.ProtocolHash
}

func NewClient(ctx context.Context, baseURL string) (*Client, error) {
	return NewClientWithHTTPClient(ctx, baseURL, nil)
}

func NewClientWithHTTPClient(ctx context.Context, baseURL string, httpClient *http.Client) (*Client, error) {
	rpcClient, err := rpc.NewClient(baseURL, httpClient)
	if err != nil {
		return nil, err
	}

	err = rpcClient.InitChain(ctx)
	if err != nil {
		return nil, err
	}

	constants, err := rpcClient.GetConstants(ctx, rpc.Head)
	if err != nil {
		return nil, err
	}

	head, err := rpcClient.GetHeadBlock(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:    rpcClient,
		Constants: constants,
		Protocol:  head.Protocol,
	}, nil
}

func (c *Client) PreApplyAndInject(ctx context.Context, signer *Signer, tx *codec.Transaction) (tezos.OpHash, error) {
	preApplyRes, op, err := c.PreApply(ctx, signer, tx)
	if err != nil {
		return tezos.ZeroOpHash, err
	}

	const (
		microTezMinFee     = 100 // µtz
		microTezPerByte    = 1   // µtz/b
		microTezPerGasUnit = 0.1 // µtz/g
		consumedGasMargin  = 100 // g
	)

	opSize := len(op.Bytes())

	computedFee := microTezMinFee + int64(microTezPerByte*opSize) + int64(microTezPerGasUnit*float64(preApplyRes.Cost().Gas+consumedGasMargin))

	return c.Inject(ctx, signer, tx, preApplyRes.Result().ConsumedGas, preApplyRes.Result().PaidStorageSizeDiff, computedFee)
}

func (c *Client) PreApply(ctx context.Context, signer *Signer, tx *codec.Transaction) (rpc.TypedOperation, *codec.Op, error) {
	op, err := c.buildOp(ctx, signer, tx, c.Constants.HardGasLimitPerOperation, c.Constants.HardStorageLimitPerOperation, 0)
	if err != nil {
		return nil, nil, err
	}

	var res []*rpc.Operation
	err = c.PreApplyOperations(ctx, rpc.Head, []PreApplyOperation{{Op: op, Protocol: c.Protocol}}, &res)
	if err != nil {
		return nil, op, err
	}

	return res[0].Contents[0], op, nil
}

func (c *Client) Inject(ctx context.Context, signer *Signer, tx *codec.Transaction, gasLimit, storageLimit, fee int64) (tezos.OpHash, error) {
	op, err := c.buildOp(ctx, signer, tx, gasLimit, storageLimit, fee)
	if err != nil {
		return tezos.ZeroOpHash, err
	}

	opHash, err := c.BroadcastOperation(ctx, op.Bytes())
	if err != nil {
		return tezos.ZeroOpHash, err
	}

	return opHash, nil
}

func (c *Client) buildOp(ctx context.Context, signer *Signer, tx *codec.Transaction, gasLimit, storageLimit, fee int64) (*codec.Op, error) {
	counter, err := c.GetCounter(ctx, signer.Address(), rpc.Head)
	if err != nil {
		return nil, err
	}

	head, err := c.GetBlockHash(ctx, rpc.Head)
	if err != nil {
		return nil, err
	}

	txMeta := &codec.Manager{Source: signer.Address()}
	txMeta = txMeta.
		WithCounter(counter + 1).
		WithGasLimit(gasLimit).
		WithStorageLimit(storageLimit).
		WithFee(fee)

	tx = &codec.Transaction{
		Manager:     *txMeta,
		Destination: tx.Destination,
		Amount:      tx.Amount,
		Parameters:  tx.Parameters,
	}

	op := codec.NewOp().
		WithBranch(head).
		WithContents(tx)

	op, err = signer.Sign(op)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func (c *Client) GetCounter(ctx context.Context, addr tezos.Address, id rpc.BlockID) (int64, error) {
	u := fmt.Sprintf("chains/main/blocks/%s/context/contracts/%s/counter", id, addr)
	var counter string
	err := c.Get(ctx, u, &counter)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(counter, 10, 64)
}

type PreApplyOperation struct {
	Protocol tezos.ProtocolHash
	*codec.Op
}

func (r *PreApplyOperation) MarshalJSON() ([]byte, error) {
	type OpAlias codec.Op
	return json.Marshal(struct {
		Protocol tezos.ProtocolHash `json:"protocol"`
		*OpAlias
	}{
		Protocol: r.Protocol,
		OpAlias:  (*OpAlias)(r.Op),
	})
}

func (c *Client) PreApplyOperations(ctx context.Context, id rpc.BlockID, body []PreApplyOperation, resp interface{}) error {
	u := fmt.Sprintf("chains/main/blocks/%s/helpers/preapply/operations", id)
	return c.Post(ctx, u, body, resp)
}
