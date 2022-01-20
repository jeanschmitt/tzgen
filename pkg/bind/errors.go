package bind

import "github.com/pkg/errors"

var (
	ErrContractNotFound = errors.New("contract not deployed")
)
