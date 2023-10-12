package goaa

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	entrypoint "github.com/pavankpdev/goaa/gen"
	factory "github.com/pavankpdev/goaa/gen"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
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
	Client      *ethclient.Client      // Ethereum client for interacting with the blockchain
	Owner       common.Address         // Ethereum address of the owner
	SAFactory   *factory.Factory       // Smart account factory contract instance
	EntryPoint  *entrypoint.EntryPoint // Smart account factory contract instance
	SignMessage func(message string) (string, error)
}

type TargetParams struct {
	Target string
	Data   string
	Value  string
}

type UOps struct {
	Sender               string  `json:"sender"`
	Nonce                string  `json:"nonce"`
	InitCode             string  `json:"initCode"`
	CallData             string  `json:"callData"`
	Signature            *string `json:"signature,omitempty"`
	CallGasLimit         string  `json:"callGasLimit"`
	VerificationGasLimit string  `json:"verificationGasLimit"`
	PreVerificationGas   string  `json:"preVerificationGas"`
	MaxFeePerGas         string  `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string  `json:"maxPriorityFeePerGas"`
	PaymasterAndData     string  `json:"paymasterAndData"`
}

type UserOperationTxnPayload struct {
	Id      int64  `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

// NewSmartAccountProvider creates a new instance of SmartAccountProvider with the provided parameters.
// It initializes the Ethereum client, owner's address, and the smart account factory contract.
func NewSmartAccountProvider(params SmartAccountProviderParams) (*SmartAccountProvider, error) {
	client, err := createEthClient(params.RPC)
	if err != nil {
		return nil, err
	}

	owner, err := privateKeyToAddress(params.OwnerPrivateKey[2:])
	if err != nil {
		return nil, err
	}

	fac, err := factory.NewFactory(common.HexToAddress(params.SmartAccountFactoryAddress), client)
	if err != nil {
		return nil, err
	}

	ep, err := entrypoint.NewEntryPoint(common.HexToAddress(params.EntryPointAddress), client)
	if err != nil {
		return nil, err
	}

	return &SmartAccountProvider{
		Client:     client,
		Owner:      owner,
		SAFactory:  fac,
		EntryPoint: ep,
		SignMessage: func(message string) (string, error) {
			signature, err := _signMessage(message, params.OwnerPrivateKey)
			if err != nil {
				return "", err
			}

			return signature, nil
		},
	}, nil
}

func _signMessage(dataToSign string, privateKeyHex string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex[2:])
	if err != nil {
		return "", err
	}

	data := []byte(dataToSign)
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signature), nil

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

func (sap *SmartAccountProvider) SendUserOpsTransaction(target TargetParams) (any, error) {
	// Nonce
	nonce, err := sap.Client.PendingNonceAt(context.Background(), sap.Owner)
	if err != nil {
		return 0, err
	}

	sender, err := sap.GetSmartAccountAddress(int64(nonce))
	if err != nil {
		return 0, err
	}

	fmt.Println(sender)

	jsonStr, err := json.Marshal(target)
	if err != nil {
		return 0, err
	}

	calldata := []byte(jsonStr)

	nonceInHex := strconv.FormatInt(int64(nonce), 16)

	uo := UOps{
		Sender:               "0xCeFc6b95c885D17E3c328f57F24b13E3EE82Aec2",
		Nonce:                "0x" + nonceInHex,
		InitCode:             "0x",
		CallData:             "0x" + hex.EncodeToString(calldata),
		CallGasLimit:         "0x2710",
		VerificationGasLimit: "0x2710",
		PreVerificationGas:   "0x402db0",
		MaxFeePerGas:         "0x17190c894e",
		MaxPriorityFeePerGas: "0x3812ed1a0",
		PaymasterAndData:     "0x",
	}

	uopsHash := crypto.Keccak256Hash([]byte(uo.Sender), []byte(uo.Nonce), []byte(uo.InitCode), []byte(uo.CallData), []byte(uo.CallGasLimit), []byte(uo.VerificationGasLimit), []byte(uo.PreVerificationGas), []byte(uo.MaxFeePerGas), []byte(uo.MaxPriorityFeePerGas), []byte(uo.PaymasterAndData))

	privateKey, err := crypto.HexToECDSA("<pvt key without 0x>")
	if err != nil {
		fmt.Println("Failed to sign the UOps struct:", err)
		return 0, err
	}

	// Fix: This Signature here is invalid
	signature, err := crypto.Sign(uopsHash.Bytes(), privateKey)
	if err != nil {
		fmt.Println("Failed to sign the UOps struct:", err)
		return 0, err
	}

	var sog = "0x" + hex.EncodeToString(signature)
	fmt.Println(sog)
	uo.Signature = &sog

	var uoArray []any
	uoArray = append(uoArray, uo)
	uoArray = append(uoArray, "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")
	fmt.Println(uoArray)

	bodyPayload, err := json.Marshal(UserOperationTxnPayload{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "eth_sendUserOperation",
		Params:  uoArray,
	})
	if err != nil {
		return 0, err
	}

	fmt.Println(string(bodyPayload))

	//res, err := sap.EntryPoint.HandleOps(nil, uoArray, sap.Owner)

	url := "https://polygon-mumbai.g.alchemy.com/v2/u-FhnHbTFL8OASxmdclXSWKS-YcypJzH"

	payload := strings.NewReader(string(bodyPayload))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body), nil

	// Create target bytes

	// generate UO object
}
