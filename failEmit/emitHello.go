package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/events"
)

var PUBLIC = sdk.Export(hello)
var SYSTEM = sdk.Export(_init)

func _init() {}

func emitHello(greeting string) {}
func hello() {
	events.EmitEvent(emitHello, "hello")
}
