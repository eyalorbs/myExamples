package tictactoe

import (
	"bytes"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
)

var PUBLIC = sdk.Export()
var SYSTEM = sdk.Export(_init)

func _init() {

}

func (game *Game) play(index uint32) (feedback string) {
	if game.PlayerXTurn {
		if bytes.Equal(address.GetCallerAddress(), game.PlayerY) {
			panic("it is player X's turn, wait for your turn")
		}
	} else {
		if bytes.Equal(address.GetCallerAddress(), game.PlayerX) {
			panic("it is player Y's turn, wait for your turn")
		}
	}

}

type Board [9]rune

type Game struct {
	Board       Board
	PlayerX     []byte
	PlayerY     []byte
	PlayerXTurn bool
}
type Games []Game
