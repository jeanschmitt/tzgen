package bind

import (
	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/tezos"
)

type Signer struct {
	privateKey tezos.PrivateKey
}

func NewSigner(privateKey string) (s *Signer, err error) {
	s = &Signer{}
	s.privateKey, err = tezos.ParsePrivateKey(privateKey)
	return s, nil
}

func (s *Signer) Address() tezos.Address {
	return s.privateKey.Address()
}

func (s *Signer) Sign(op *codec.Op) (*codec.Op, error) {
	err := op.Sign(s.privateKey)
	if err != nil {
		return nil, err
	}
	return op, err
}
