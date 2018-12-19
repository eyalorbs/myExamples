package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/safemath/safeuint64"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
)

var PUBLIC = sdk.Export(totalSupply, balanceOf, allowance, transfer, approve, transferFrom, 	getSymbol, getName, getDecimals)
var SYSTEM = sdk.Export(_init)

func _init(){

	state.WriteStringByKey("symbol", "E20")
	state.WriteStringByKey("name", "ERC20Token")
	state.WriteUint32ByKey("decimals", 18)
	state.WriteUint64ByKey("totalSupply", 1000000000000000000)
	state.WriteUint64ByAddress(address.GetSignerAddress(), totalSupply())
}

//the following functions are not an ERC20 requirement, but it is necessary in order to read these variables from the state
func getSymbol() string{
	return state.ReadStringByKey("symbol")
}

func getName() string{
	return state.ReadStringByKey("name")
}

func getDecimals()uint32{
	return state.ReadUint32ByKey("decimals")
}






func totalSupply() (amount uint64){
	return state.ReadUint64ByKey("totalSupply")
}

func balanceOf(tokenOwner []byte)(balance uint64){
	//validate the address
	address.ValidateAddress(tokenOwner)

	return state.ReadUint64ByAddress(tokenOwner)
}

//we will declare that the key that will indicate the allowance of the spender is a
//concatenation of the tokenOwner address and the spender address
//we do not use maps because as of writing this the orbs network doesn't support maps


func allowance(tokenOwner, spender []byte) (remaining uint64){

	//ensure that the addresses are valid
	address.ValidateAddress(tokenOwner)
	address.ValidateAddress(spender)

	//get the key
	key := append(tokenOwner, spender...)

	return state.ReadUint64ByAddress(key)
}




func transfer(to []byte, tokens uint64){

	//validate the address
	address.ValidateAddress(to)

	//update spender's tokens
	state.WriteUint64ByAddress(address.GetCallerAddress(), safeuint64.Sub(state.ReadUint64ByAddress(address.GetCallerAddress()), tokens))

	//update receiver's tokens
	state.WriteUint64ByAddress(to, safeuint64.Add(state.ReadUint64ByAddress(to), tokens))

	//TODO: add an event

}

func approve(spender []byte, tokens uint64){
	//get the key
	key := append(address.GetCallerAddress(), spender...)

	state.WriteUint64ByAddress(key, safeuint64.Add(state.ReadUint64ByAddress(key), tokens))

	//TODO: add an event
}

func transferFrom(from, to []byte, tokens uint64){
	//get the key
	key := append(from, address.GetCallerAddress()...)

	//update token owner account
	state.WriteUint64ByAddress(from, safeuint64.Sub(state.ReadUint64ByAddress(from), tokens))

	//update allowance
	state.WriteUint64ByAddress(key, safeuint64.Sub(state.ReadUint64ByAddress(key), tokens))

	//update receiver's tokens
	state.WriteUint64ByAddress(to, safeuint64.Add(state.ReadUint64ByAddress(key), tokens))
}

