package goaa

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	factory "github.com/pavankpdev/goaa/gen"
	"math/big"
)

// SmartAccountProviderParams stores the parameters required to initialize the SmartAccountProvider.
type SmartAccountProviderParams struct {
	OwnerPrivateKey            string // The private key of the Ethereum account
	RPC                        string // The RPC endpoint for the Ethereum node
	EntryPointAddress          string // The address of the entry point contract
	SmartAccountFactoryAddress string // The address of the smart account factory contract
}

// SmartAccountProvider is a struct that manages interaction with Ethereum smart contracts.
type SmartAccountProvider struct {
	Client    *ethclient.Client // Ethereum client for interacting with the blockchain
	Owner     common.Address    // Ethereum address of the owner
	SAFactory *factory.Factory  // Smart account factory contract instance
}

// NewSmartAccountProvider creates a new instance of SmartAccountProvider with the provided parameters.
// It initializes the Ethereum client, owner's address, and the smart account factory contract.
func NewSmartAccountProvider(params SmartAccountProviderParams) (*SmartAccountProvider, error) {
	client, err := createEthClient(params.RPC)
	if err != nil {
		return nil, err
	}

	owner, err := privateKeyToAddress(params.OwnerPrivateKey)
	if err != nil {
		return nil, err
	}

	fac, err := factory.NewFactory(common.HexToAddress(params.SmartAccountFactoryAddress), client)
	if err != nil {
		return nil, err
	}

	return &SmartAccountProvider{
		Client:    client,
		Owner:     owner,
		SAFactory: fac,
	}, nil
}

// createEthClient connects to an Ethereum node via the specified RPC endpoint
// and returns an Ethereum client. It panics on connection errors.
func createEthClient(rpc string) (*ethclient.Client, error) {
	cl, err := ethclient.Dial(rpc)

	if err != nil {
		return nil, err
	}

	return cl, nil
}

// privateKeyToAddress converts a private key (in hexadecimal format) to an Ethereum address.
func privateKeyToAddress(privateKey string) (common.Address, error) {
	key, err := crypto.HexToECDSA(privateKey)

	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(key.PublicKey), nil
}

// GetSmartAccountAddress retrieves the address of a smart account based on a given salt value.
func (sap *SmartAccountProvider) GetSmartAccountAddress(salt int64) (common.Address, error) {

	address, err := sap.SAFactory.GetAddress(nil, sap.Owner, big.NewInt(salt))
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}
