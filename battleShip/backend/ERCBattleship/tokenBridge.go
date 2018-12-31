package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/service"
)

var PUBLIC = sdk.Export(transferToContract)
var SYSTEM = sdk.Export(_init)

func _init() {

}

func transferToContract(from []byte, tokens uint64) {
	service.CallMethod("token", "transferFrom", from, address.GetCallerAddress(), tokens)
}
