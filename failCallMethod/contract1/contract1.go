package main

import "github.com/orbs-network/orbs-contract-sdk/go/sdk"

var PUBLIC = sdk.Export(addOne)
var SYSTEM = sdk.Export(_init)

func _init() {

}

func addOne(x uint32) (sum uint32) {
	return x + 1
}
