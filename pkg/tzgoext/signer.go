package tzgoext

import (
	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/tezos"
	"context"
)

type Signer struct {
	key tezos.PrivateKey
}

func NewSigner(key tezos.PrivateKey) *Signer {
	return &Signer{key: key}
}

func NewSignerFromPrivateKeyString(s string) (*Signer, error) {
	key, err := tezos.ParsePrivateKey(s)
	if err != nil {
		return nil, err
	}
	return NewSigner(key), nil
}

func (s *Signer) Address(_ context.Context) (tezos.Address, error) {
	return s.key.Address(), nil
}

func (s *Signer) Key(_ context.Context) (tezos.Key, error) {
	return s.key.Public(), nil
}

func (s *Signer) SignMessage(_ context.Context, s2 string) (tezos.Signature, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Signer) SignOperation(_ context.Context, op *codec.Op) (tezos.Signature, error) {
	return s.key.Sign(op.Digest())
}

func (s *Signer) SignBlock(_ context.Context, header *codec.BlockHeader) (tezos.Signature, error) {
	sig, err := s.key.Sign(header.Digest())
	if err != nil {
		return tezos.Signature{}, err
	}
	sig.Type = tezos.SignatureTypeGeneric
	return sig, nil
}
