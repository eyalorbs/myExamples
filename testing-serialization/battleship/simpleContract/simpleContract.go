package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/events"
)

var PUBLIC = sdk.Export(returnAndEvent, onlyReturn, onlyEvent)
var SYSTEM = sdk.Export(_init, event1)

func _init() {

}
func event1(greeting string) {
}

func returnAndEvent() (firstGreeting string) {
	events.EmitEvent(event1, "hello good sir,\n how are you")
	return "hello here is a new line\n and here I am after a new line"
}

func onlyReturn() (greeting string) {
	return "hello\nhello good day"
}

func onlyEvent() {
	events.EmitEvent(event1, "if you see me it is because there is a god that let you see this")
}
