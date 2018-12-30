package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/service"
)

var PUBLIC = sdk.Export(transferToContract, transferToUser, getBattleShipBalance, getUserBalance)
var SYSTEM = sdk.Export(_init)

func _init() {

}

func transferToContract(tokens uint64) {
	service.CallMethod("token", "transfer", address.GetCallerAddress(), tokens)
}

func transferToUser(tokens uint64) {
	service.CallMethod("token", "approve", 7)
	service.CallMethod("token", "transfer", address.GetSignerAddress(), tokens)
}
func getBattleShipBalance() (tokens uint64) {
	value := service.CallMethod("token", "balanceOf", address.GetCallerAddress())[0]
	if tokens, ok := value.(uint64); ok {
		return tokens
	}
	panic("invalid return argument")
}

func getUserBalance() (tokens uint64) {
	value := service.CallMethod("token", "balanceOf", address.GetSignerAddress())[0]
	if tokens, ok := value.(uint64); ok {
		return tokens
	}
	panic("invalid return argument")
}
