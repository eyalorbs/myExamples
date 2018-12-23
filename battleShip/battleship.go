package main

import (
	"bytes"
	"encoding/json"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/events"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
)

var PUBLIC = sdk.Export(startNewGame, guess)
var SYSTEM = sdk.Export(_init)

func _init() {
	state.WriteBytesByKey("waitingPool", []byte{})

	var games = Games{}
	b, err := games.MarshalJSON()
	if err != nil {
		panic(err)
	}
	state.WriteBytesByKey("games", b)

}

func startNewGame(merkleRoot []byte) {
	signerAddress := address.GetSignerAddress()

	//check if player is already playing
	if !bytes.Equal(state.ReadBytesByAddress(address.GetSignerAddress()), []byte{}) {
		panic("you are already playing a game")
	}

	//check if there is someone in the waiting pool
	if len(state.ReadBytesByKey("waitingPool")) == 0 {
		state.WriteBytesByKey("waitingPool", signerAddress)
		state.WriteBytesByAddress(signerAddress, merkleRoot)
		return
	}
	//get player2 address
	player2 := state.ReadBytesByKey("waitingPool")[:20]

	//update the waiting pool
	state.WriteBytesByKey("waitingPool", state.ReadBytesByKey("waitingPool")[20:])

	//create new game
	var game Game
	game.new(address.GetSignerAddress(), state.ReadBytesByKey("games"), merkleRoot, state.ReadBytesByAddress(player2))

	//get games from state
	games := Games{}
	games.getGames()
	games[string(append(signerAddress, player2...))] = game

	b, err := games.MarshalJSON()
	if err != nil {
		panic(err)
	}
	state.WriteBytesByKey("games", b)
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////todo understand events
	// todo understand events
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	events.EmitEvent(startNewGame, signerAddress, player2)
}

func guess(x, y uint32) {
	events.EmitEvent(guess, x, y, address.GetSignerAddress())
}

type Game struct {
	Player1     []byte
	Player2     []byte
	Player1Root []byte
	Player2Root []byte
}

func (game *Game) new(player1, player2, player1root, player2root []byte) {
	game.Player1 = player1
	game.Player2 = player2
	game.Player1Root = player1root
	game.Player2Root = player2root
}
func (game *Game) MarshalJSON() (b []byte, err error) {
	gameMap := make(map[byte][]byte)
	gameMap[0] = game.Player1
	gameMap[1] = game.Player2
	gameMap[2] = game.Player1Root
	gameMap[3] = game.Player2Root
	return json.Marshal(gameMap)
}
func (game *Game) UnmarshalJSON(b []byte) (err error) {
	gameMap := make(map[byte][]byte)
	err = json.Unmarshal(b, &gameMap)
	game.Player1 = gameMap[0]
	game.Player2 = gameMap[1]
	game.Player1Root = gameMap[2]
	game.Player2Root = gameMap[3]
	return nil
}

type Games map[string]Game

func (games *Games) MarshalJSON() (b []byte, err error) {
	gamesMap := make(map[string][]byte)
	for key, value := range *games {
		gamesMap[key], err = value.MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
	}
	return json.Marshal(gamesMap)
}
func (games *Games) UnmarshalJSON(b []byte) (err error) {
	var gamesMap = make(map[string][]byte)
	err = json.Unmarshal(b, &gamesMap)
	if err != nil {
		return err
	}
	var tempGame Game
	for key, value := range gamesMap {
		err = json.Unmarshal(value, &tempGame)
		if err != nil {
			return err
		}
		(*games)[key] = tempGame
	}
	return nil
}
func (games *Games) getGames() {
	gamesBytes := state.ReadBytesByKey("games")
	err := games.UnmarshalJSON(gamesBytes)
	if err != nil {
		panic(err)
	}
}
