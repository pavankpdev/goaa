package utils

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	goaa "github.com/pavankpdev/goaa"
)

func BuildUserOp(sender common.Address, nonce string, calldata []byte) goaa.UOps {
	return goaa.UOps{
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
