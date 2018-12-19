package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/service"
)

var PUBLIC = sdk.Export(add, get)
var SYSTEM = sdk.Export(_init)

func _init(){

}

func add(amount uint64){
	service.CallMethod("MyCounter", "add", amount)
}

func get()(count uint64){
	temp := service.CallMethod("MyCounter", "get")
	if num, ok :=temp[0].(uint64); ok{
		return num
	}
	panic("unexpected return from counter contract")
}