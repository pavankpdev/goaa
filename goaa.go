package goaa

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

// CreateClient connects to an Ethereum node via the specified RPC endpoint
// and returns an Ethereum client. It panics on connection errors.
func CreateClient(rpc string) *ethclient.Client {
	cl, err := ethclient.Dial(rpc)

	if err != nil {
		panic(err)
	}

	return cl
}
