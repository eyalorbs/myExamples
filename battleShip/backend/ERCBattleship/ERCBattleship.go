package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/service"
)

var PUBLIC = sdk.Export(transfer, getBattleShipBalance, getUserBalance)
var SYSTEM = sdk.Export(_init)

func _init() {

}

func transfer(tokens uint64) {
	service.CallMethod("ERC20Token", "transfer", address.GetCallerAddress(), tokens)
}

func getBattleShipBalance() (tokens uint64) {
	value := service.CallMethod("ERC20Token", "balanceOf", address.GetCallerAddress())[0]
	if tokens, ok := value.(uint64); ok {
		return tokens
	}
	panic("invalid return argument")
}

func getUserBalance() (tokens uint64) {
	value := service.CallMethod("ERC20Token", "balanceOf", address.GetSignerAddress())[0]
	if tokens, ok := value.(uint64); ok {
		return tokens
	}
	panic("invalid return argument")
}
