package tictactoe

import (
	"bytes"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
)

var PUBLIC = sdk.Export()
var SYSTEM = sdk.Export(_init)

func _init() {
	state.WriteBytesByKey("waitingPool", []byte{})
	state.WriteBytesByKey("games", []byte{})
}

func startGame() {
	callerAddress := address.GetCallerAddress()

	if len(state.ReadBytesByKey("waitingPool")) == 0 {
		state.WriteBytesByKey("waitingPool", callerAddress)
	} else {

		state.WriteUint32ByAddress(callerAddress)

	}
}
func (game *Game) play(index uint32) {

	callerAddress := address.GetCallerAddress()
	if !(bytes.Equal(callerAddress, game.PlayerX) || bytes.Equal(callerAddress, game.PlayerO)) {
		panic("you are not registered as a player in this game")
	}

	if game.PlayerXTurn {
		if bytes.Equal(callerAddress, game.PlayerO) {
			panic("it is player X's turn, wait for your turn")
		}
	} else {
		if bytes.Equal(callerAddress, game.PlayerX) {
			panic("it is player Y's turn, wait for your turn")
		}
	}

	if 8 < index {
		panic("the index must be between 0 and 8")
	}

	if game.Board[index] != '-' {
		panic("that index is taken")
	}

	if game.PlayerXTurn {
		game.Board[index] = 'X'
	} else {
		game.Board[index] = 'Y'
	}

}

type Board [9]rune

type Game struct {
	Board       Board
	PlayerX     []byte
	PlayerO     []byte
	PlayerXTurn bool
}
type Games []Game
