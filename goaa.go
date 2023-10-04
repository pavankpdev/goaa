package goaa

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

// PrivateKeyToAddress converts a private key (in hexadecimal format) to an Ethereum address.
func PrivateKeyToAddress(pvtKey string) common.Address {
	privateKey, err := crypto.HexToECDSA(pvtKey)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()

	// Check if the public key can be type-asserted to *ecdsa.PublicKey.
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	// Convert the ECDSA public key to an Ethereum address.
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address
}
