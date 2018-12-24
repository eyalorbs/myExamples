package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/service"
)

var PUBLIC = sdk.Export(addOne)
var SYSTEM = sdk.Export(_init)

func _init() {

}
func addOne() (number uint32) {
	value := service.CallMethod("contract1", "addOne", uint32(5))[0] // when not casting, an error appears
	if number, ok := value.(uint32); ok {
		return number
	}
	panic("invalid response")
}
