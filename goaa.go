package goaa

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	contracts := &ContractAddressParams{
		factory:    params.SmartAccountFactoryAddress,
		entrypoint: params.EntryPointAddress,
	}

	return &SmartAccountProvider{
		Client:     client,
		Owner:      owner,
		SAFactory:  fac,
		EntryPoint: ep,
		PrivateKey: params.OwnerPrivateKey,
		Contracts:  contracts,
	}, nil
}

func (sap *SmartAccountProvider) signMessage(dataToSign []byte, privateKey *ecdsa.PrivateKey) (common.Hash, error) {
	nonce, err := sap.Client.PendingNonceAt(context.Background(), sap.Owner)
	if err != nil {
		return common.Hash{}, err
	}

	to := common.HexToAddress(sap.Contracts.entrypoint)
	tx := types.NewTx(&types.AccessListTx{
		Nonce:    nonce,
		GasPrice: big.NewInt(20000000000),
		Gas:      uint64(21000),
		To:       &to,
		Value:    big.NewInt(1000000000000000000),
		Data:     dataToSign,
	})

	signature, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(80001)), privateKey)

	if err != nil {
		return common.Hash{}, err
	}

	return signature.Hash(), nil

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

func buildUserOp(sender common.Address, nonce string, calldata []byte) UOps {
	return UOps{
		Sender:               sender,
		Nonce:                "0x" + nonce,
		InitCode:             "0x",
		CallData:             "0x" + hex.EncodeToString(calldata),
		CallGasLimit:         "0x2710",
		VerificationGasLimit: "0x2710",
		PreVerificationGas:   "0x402db0",
		MaxFeePerGas:         "0x17190c894e",
		MaxPriorityFeePerGas: "0x3812ed1a0",
		PaymasterAndData:     "0x",
	}
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

	calldata, err := json.Marshal(target)
	if err != nil {
		return 0, err
	}

	nonceInHex := strconv.FormatInt(int64(nonce), 16)

	uo := buildUserOp(sender, nonceInHex, calldata)

	privateKey, err := crypto.HexToECDSA(sap.PrivateKey)
	if err != nil {
		fmt.Println("Failed to sign the UOps struct:", err)
		return 0, err
	}

	signature, err := sap.signMessage(calldata, privateKey)

	if err != nil {
		fmt.Println("Failed to sign the UOps struct:", err)
		return 0, err
	}

	uo.Signature = &signature

	var uoArray []any
	uoArray = append(uoArray, uo)
	uoArray = append(uoArray, sap.Contracts.entrypoint)

	bodyPayload, err := json.Marshal(UserOperationTxnPayload{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "eth_sendUserOperation",
		Params:  uoArray,
	})

	if err != nil {
		return 0, err
	}

	url := "https://polygon-mumbai.g.alchemy.com/v2/u-FhnHbTFL8OASxmdclXSWKS-YcypJzH"

	payload := strings.NewReader(string(bodyPayload))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	return string(body), nil

}
