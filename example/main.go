package main

import (
	"fmt"
	"github.com/pavankpdev/goaa"
)

func main() {

	const RPC = "<RPC>"
	const SmartAccountFactoryAddress = "0x9406Cc6185a346906296840746125a0E44976454" // https://docs.alchemy.com/docs/creating-a-smart-contract-account-and-sending-userops#1b-get-constants
	const PrivateKey = "<PrivateKey>"

	params := goaa.SmartAccountProviderParams{
		OwnerPrivateKey:            PrivateKey,
		RPC:                        RPC,
		EntryPointAddress:          "",
		SmartAccountFactoryAddress: SmartAccountFactoryAddress,
	}

	client, err := goaa.NewSmartAccountProvider(params)

	if err != nil {
		panic(err)
	}

	address, err := client.GetSmartAccountAddress(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("My samrt account address is %v\n", address)

}
