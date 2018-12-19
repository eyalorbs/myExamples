package main

import (
"github.com/orbs-network/orbs-contract-sdk/go/sdk"
"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
)

var PUBLIC = sdk.Export(greet)
var SYSTEM = sdk.Export(_init)

func _init() {
state.WriteStringByKey("greeting", "hello world!")
}



func greet() string {
	return state.ReadStringByKey("greeting")
}