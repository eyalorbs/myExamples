package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
	"math"
)

var PUBLIC = sdk.Export(totalSupply, balanceOf, allowance, transfer, approve, transferFrom, getSymbol, getName, getDecimals)
var SYSTEM = sdk.Export(_init)

func _init() {

	state.WriteStringByKey("symbol", "E20")
	state.WriteStringByKey("name", "ERC20Token")
	state.WriteUint32ByKey("decimals", 18)
	state.WriteUint64ByKey("totalSupply", 1000000000000000000)
	state.WriteUint64ByAddress(address.GetSignerAddress(), totalSupply())
}

//the following functions are not an ERC20 requirement, but it is necessary in order to read these variables from the state
func getSymbol() string {
	return state.ReadStringByKey("symbol")
}

func getName() string {
	return state.ReadStringByKey("name")
}

func getDecimals() uint32 {
	return state.ReadUint32ByKey("decimals")
}

func totalSupply() (amount uint64) {
	return state.ReadUint64ByKey("totalSupply")
}

func balanceOf(tokenOwner []byte) (balance uint64) {
	//validate the address
	address.ValidateAddress(tokenOwner)

	return state.ReadUint64ByAddress(tokenOwner)
}

//we will declare that the key that will indicate the allowance of the spender is a
//concatenation of the tokenOwner address and the spender address
//we do not use maps because as of writing this the orbs network doesn't support maps

func allowance(tokenOwner, spender []byte) (remaining uint64) {

	//ensure that the addresses are valid
	address.ValidateAddress(tokenOwner)
	address.ValidateAddress(spender)

	//get the key
	key := append(tokenOwner, spender...)

	return state.ReadUint64ByAddress(key)
}

func transfer(to []byte, tokens uint64) {

	//validate the address
	address.ValidateAddress(to)

	//update spender's tokens
	state.WriteUint64ByAddress(address.GetSignerAddress(), Sub(state.ReadUint64ByAddress(address.GetSignerAddress()), tokens))

	//update receiver's tokens
	state.WriteUint64ByAddress(to, Add(state.ReadUint64ByAddress(to), tokens))

	//TODO: add an event

}

func approve(spender []byte, tokens uint64) {
	//validate address
	address.ValidateAddress(spender)
	//get the key
	key := append(address.GetSignerAddress(), spender...)
	state.WriteUint64ByAddress(key, Add(state.ReadUint64ByAddress(key), tokens))

	//TODO: add an event
}

func transferFrom(from, to []byte, tokens uint64) {
	//validate addresses
	address.ValidateAddress(from)
	address.ValidateAddress(to)

	//get the key
	s := []byte{}
	s = append(s, from...)
	s = append(s, address.GetCallerAddress()...)
	key := s

	//update token owner account
	state.WriteUint64ByAddress(from, Sub(state.ReadUint64ByAddress(from), tokens))

	//update allowance
	state.WriteUint64ByAddress(key, Sub(state.ReadUint64ByAddress(key), tokens))

	//update receiver's tokens
	state.WriteUint64ByAddress(to, Add(state.ReadUint64ByAddress(to), tokens))
}

//delete these functions once gamma is updated:

func Add(x uint64, y uint64) uint64 {
	if y > math.MaxUint64-x {
		panic("integer overflow on add")
	}
	return x + y
}

func Sub(x uint64, y uint64) uint64 {
	if x < y {
		panic("integer overflow on sub")
	}
	return x - y
}

func Mul(x uint64, y uint64) uint64 {
	if x == 0 || y == 0 {
		return 0
	}
	if y > math.MaxUint64/x {
		panic("integer overflow on mul")
	}
	return x * y
}

func Mod(x uint64, y uint64) uint64 {
	if y == 0 {
		panic("division by zero")
	}
	return x % y
}
