package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
)

var PUBLIC = sdk.Export(writeWinner, getWinner)
var SYSTEM = sdk.Export(_init)

func _init() {}

func writeWinner(player []byte, didWin uint32) {
	state.WriteUint32ByAddress(player, didWin)
}

func getWinner(player []byte) (didWin uint32) {
	return state.ReadUint32ByAddress(player)
}
