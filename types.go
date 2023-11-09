package goaa

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	entrypoint "github.com/pavankpdev/goaa/gen"
	factory "github.com/pavankpdev/goaa/gen"
)

// SmartAccountProviderParams stores the parameters required to initialize the SmartAccountProvider.
type SmartAccountProviderParams struct {
	OwnerPrivateKey            string // The private key of the Ethereum account
	RPC                        string // The RPC endpoint for the Ethereum node
	EntryPointAddress          string // The address of the entry point contract
	SmartAccountFactoryAddress string // The address of the smart account factory contract
}

type ContractAddressParams struct {
	factory    string
	entrypoint string
}

// SmartAccountProvider is a struct that manages interaction with Ethereum smart contracts.
type SmartAccountProvider struct {
	Client     *ethclient.Client      // Ethereum client for interacting with the blockchain
	Owner      common.Address         // Ethereum address of the owner
	SAFactory  *factory.Factory       // Smart account factory contract instance
	EntryPoint *entrypoint.EntryPoint // Smart account factory contract instance
	PrivateKey string                 // The private key of the Ethereum account
	Contracts  *ContractAddressParams // The object that contains all the contract addresses
}

type TargetParams struct {
	Target string
	Data   string
	Value  string
}

type UOps struct {
	Sender               common.Address `json:"sender"`
	Nonce                string         `json:"nonce"`
	InitCode             string         `json:"initCode"`
	CallData             string         `json:"callData"`
	Signature            *common.Hash   `json:"signature,omitempty"`
	CallGasLimit         string         `json:"callGasLimit"`
	VerificationGasLimit string         `json:"verificationGasLimit"`
	PreVerificationGas   string         `json:"preVerificationGas"`
	MaxFeePerGas         string         `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string         `json:"maxPriorityFeePerGas"`
	PaymasterAndData     string         `json:"paymasterAndData"`
}

type UserOperationTxnPayload struct {
	Id      int64  `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}
