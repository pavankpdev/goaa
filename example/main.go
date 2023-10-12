package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"github.com/pavankpdev/goaa"
	"math/big"
)

func main() {

	const RPC = "https://polygon-mumbai.g.alchemy.com/v2/u-FhnHbTFL8OASxmdclXSWKS-YcypJzH"
	const SmartAccountFactoryAddress = "0x9406Cc6185a346906296840746125a0E44976454" // https://docs.alchemy.com/docs/creating-a-smart-contract-account-and-sending-userops#1b-get-constants
	const PrivateKey = "<pvt key>"

	eth := big.NewFloat(0.1)
	wei := new(big.Float)
	wei.Mul(eth, big.NewFloat(params.Ether))

	SAParams := goaa.SmartAccountProviderParams{
		OwnerPrivateKey:            PrivateKey,
		RPC:                        RPC,
		EntryPointAddress:          "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789",
		SmartAccountFactoryAddress: SmartAccountFactoryAddress,
	}

	client, err := goaa.NewSmartAccountProvider(SAParams)

	if err != nil {
		fmt.Println("Error creating SmartAccountProvider:", err)
		panic(err)
	}

	address, err := client.GetSmartAccountAddress(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("My samrt account address is %v\n", address)

	nonce, err := client.SendUserOpsTransaction(goaa.TargetParams{
		Target: "0x94f3178AcB40d0E9c6967108e3711CF047D3240A",
		Data:   "0x",
		Value:  wei.String(),
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("My samrt account address is %v\n", nonce)

}
