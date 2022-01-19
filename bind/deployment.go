package bind

import "github.com/completium/go-tezos/v4/rpc"

func IsDeployed(tzClient *rpc.Client, address string) (bool, error) {
	res, err := tzClient.ContractStorage(rpc.ContractStorageInput{
		ContractID: address,
		BlockID:    &rpc.BlockIDHead{},
	})
	if err != nil {
		return false, err
	}
	return res.RawResponse.StatusCode == 200, nil
}
