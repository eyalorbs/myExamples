package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/service"
)

var PUBLIC = sdk.Export(transferToContract, transferToUser, getContractBalance, getPlayerBalance)
var SYSTEM = sdk.Export(_init)

func _init() {

}

func transferToContract() {
	service.CallMethod("ERCBattleship", "transferToContract", 10)
}

func transferToUser() {
	service.CallMethod("ERCBattleship", "transferToContract", 5)

}

func getContractBalance() (tokens uint64) {
	return service.CallMethod("ERCBattleship", "getBattleShipBalance")[0].(uint64)
}
func getPlayerBalance() (tokens uint64) {
	return service.CallMethod("ERCBattleship", "getUserBalance")[0].(uint64)
}
