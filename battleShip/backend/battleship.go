package main

import (
	"bytes"
	"encoding/json"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/service"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
	"math"
)

var PUBLIC = sdk.Export(startGame, getContractBalance, getUserBalance, getOpponentStatus)
var SYSTEM = sdk.Export(_init, PanicIfSignerPlaying, PanicIfSignerNotPlaying)

//helper coin contract name: ERCBattleship

func _init() {
	games := make(games)
	b, err := games.MarshalJSON()
	if err != nil {
		panic(err)
	}
	state.WriteBytesByKey("games", b)
	state.WriteBytesByKey("waitingPool", []byte{})
}

func PanicIfSignerPlaying() {
	signerAddress := address.GetSignerAddress()

	//if there are 256 bits(a hashedBoard) than the player is playing
	if len(state.ReadBytesByAddress(signerAddress)) == 32 {
		panic("player is already in pool")
	}

	//if the index isn't 0, the player is playing
	index := state.ReadUint64ByAddress(signerAddress)
	if index != 0 {
		panic("player is already playing")
	}

}

func PanicIfSignerNotPlaying() {
	signerAddress := address.GetSignerAddress()

	//if there are 256 bits(a hashedBoard) than the player is playing
	if len(state.ReadBytesByAddress(signerAddress)) == 32 {
		panic("player is in pool")
	}

	//if the index isn't 0, the player is playing
	index := state.ReadUint64ByAddress(signerAddress)
	if index == 0 {
		panic("player is not playing already playing")
	}
}

//public funtions
func startGame(hashedBoard []byte) {
	//get the games from the state
	games := make(games)
	games.getGamesFromState()
	//make sure the player isn't playing
	PanicIfSignerPlaying()

	//if there isn't anyone in the pool, add the player to the pool and save his hash map. plus he pays
	if bytes.Equal(state.ReadBytesByKey("waitingPool"), []byte{}) {
		state.WriteBytesByKey("waitingPool", address.GetSignerAddress())
		state.WriteBytesByAddress(address.GetSignerAddress(), hashedBoard)
		service.CallMethod("ERCBattleship", "transfer", 10)
		return
	}
	//player1 is the caller address
	player1 := address.GetSignerAddress()

	//get player2Address from the pool
	player2 := state.ReadBytesByKey("waitingPool")[:20]
	//update the waiting pool
	state.WriteBytesByKey("waitingPool", state.ReadBytesByKey("waitingPool")[20:])

	//get the values necessary to start a new game
	board1Hashed := hashedBoard
	//player2 hashed board is read from the state
	board2Hashed := state.ReadBytesByAddress(player2)
	defaultPlayerHits := uint8(0)
	defaultPlayerGuesses := guesses{}
	turnDefault := true
	defaultLastGuess := coordinate{}
	var newGame game
	newGame.new(player1, player2, board1Hashed, board2Hashed, defaultPlayerHits,
		defaultPlayerHits, defaultPlayerGuesses, defaultPlayerGuesses, turnDefault, defaultLastGuess)

	//check which index is free
	for i := uint64(1); i < math.MaxUint64; i++ {
		if _, ok := games[i]; !ok {
			//add the newGame go the games
			games[i] = newGame
			//get the marshaled games and add to the state
			b, err := games.MarshalJSON()
			if err != nil {
				panic(err)
			}
			state.WriteBytesByKey("games", b)

			//update the state: player's address are read as the index of the game they are playing
			state.WriteUint64ByAddress(address.GetSignerAddress(), i)
			state.ClearByAddress(player2)
			state.WriteUint64ByAddress(player2, i)
			service.CallMethod("ERCBattleship", "transfer", 10)
			return
		}
	}
	panic("no more room for more games, this game is very very very successful")
}

func getOpponentStatus() (x, y uint32) {
	signerAddress := address.GetSignerAddress()
	PanicIfSignerNotPlaying()
	//get games and index of player
	index := state.ReadUint64ByAddress(signerAddress)
	games := make(games)
	err := games.UnmarshalJSON(state.ReadBytesByKey("games"))
	if err != nil {
		panic(err)
	}
	//if the opponent hasn't played yed, panic. otherwise return the opponent guess
	opponentGuess := games[index].OpponentLastGuess
	EmptyGuess := coordinate{}
	if opponentGuess != EmptyGuess {
		panic("opponent didn't play yet")
	}
	return uint32(opponentGuess.X), uint32(opponentGuess.Y)
}

func getContractBalance() (tokens uint64) {
	value := service.CallMethod("ERCBattleship", "getBattleShipBalance")[0]
	if tokens, ok := value.(uint64); ok {
		return tokens
	}
	panic("invalid return value")
}

func getUserBalance() (tokens uint64) {
	value := service.CallMethod("ERCBattleship", "getUserBalance")[0]
	if tokens, ok := value.(uint64); ok {
		return tokens
	}
	panic("invalid return value")
}

type coordinate struct {
	X uint8
	Y uint8
}

func (coo *coordinate) new(x, y uint8) {
	coo.X = x
	coo.Y = y
}
func (coo *coordinate) MarshalJSON() (b []byte, err error) {
	coordinateMap := make(map[rune]uint8)
	coordinateMap['X'] = coo.X
	coordinateMap['Y'] = coo.Y
	return json.Marshal(coordinateMap)
}
func (coo *coordinate) UnmarshalJSON(b []byte) (err error) {
	coordinateMap := make(map[rune]uint8)
	err = json.Unmarshal(b, &coordinateMap)
	if err != nil {
		return err
	}
	coo.X = coordinateMap['X']
	coo.Y = coordinateMap['Y']
	return nil
}

type guesses []coordinate

func (guesses *guesses) MarshalJSON() (b []byte, err error) {
	length := uint8(len(*guesses))
	guessesMap := make(map[uint8][]byte)
	guessesMap[0], err = json.Marshal(length)
	if err != nil {
		return []byte{}, err
	}
	for i := uint8(0); i < length; i++ {
		guessesMap[i+1], err = (*guesses)[i].MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
	}
	return json.Marshal(guessesMap)
}
func (guesses *guesses) UnmarshalJSON(b []byte) (err error) {
	guessesMap := make(map[uint8][]byte)
	err = json.Unmarshal(b, &guessesMap)
	if err != nil {
		return err
	}
	var length uint8
	err = json.Unmarshal(guessesMap[0], &length)
	if err != nil {
		return err
	}
	for i := uint8(0); i < length; i++ {
		var tempGuess coordinate
		err = tempGuess.UnmarshalJSON(guessesMap[i+1])
		if err != nil {
			return err
		}
		*guesses = append(*guesses, tempGuess)
	}
	return nil

}

type game struct {
	Player1           []byte
	Player2           []byte
	Board1Hashed      []byte
	Board2Hashed      []byte
	Player1Hits       uint8
	Player2Hits       uint8
	Player1Guesses    guesses
	Player2Guesses    guesses
	Player1Turn       bool
	OpponentLastGuess coordinate
}

func (game *game) new(player1, player2, board1Hashed, board2Hashed []byte, player1Hits, player2Hits uint8,
	player1Guesses, player2Guesses guesses, player1Turn bool, opponentLastGuess coordinate) {
	game.Player1 = player1
	game.Player2 = player2
	game.Board1Hashed = board1Hashed
	game.Board2Hashed = board2Hashed
	game.Player1Hits = player1Hits
	game.Player2Hits = player2Hits
	game.Player1Guesses = player1Guesses
	game.Player2Guesses = player2Guesses
	game.Player1Turn = player1Turn
	game.OpponentLastGuess = opponentLastGuess
}
func (game *game) MarshalJSON() (b []byte, err error) {
	gameMap := make(map[uint8][]byte)
	gameMap[0] = game.Player1
	gameMap[1] = game.Player2
	gameMap[2] = game.Board1Hashed
	gameMap[3] = game.Board2Hashed
	gameMap[4] = []byte{game.Player1Hits}
	gameMap[5] = []byte{game.Player2Hits}
	gameMap[6], err = game.Player1Guesses.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}
	gameMap[7], err = game.Player2Guesses.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}
	if game.Player1Turn {
		gameMap[8] = []byte{1}
	} else {
		gameMap[8] = []byte{0}
	}
	gameMap[9], err = game.OpponentLastGuess.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(gameMap)

}
func (game *game) UnmarshalJSON(b []byte) (err error) {
	gameMap := make(map[uint8][]byte)
	err = json.Unmarshal(b, &gameMap)
	if err != nil {
		return err
	}

	game.Player1 = gameMap[0]
	game.Player2 = gameMap[1]
	game.Board1Hashed = gameMap[2]
	game.Board2Hashed = gameMap[3]
	game.Player1Hits = gameMap[4][0]
	game.Player2Hits = gameMap[5][0]

	err = game.Player1Guesses.UnmarshalJSON(gameMap[6])
	if err != nil {
		return err
	}

	err = game.Player2Guesses.UnmarshalJSON(gameMap[7])
	if err != nil {
		return err
	}
	if gameMap[8][0] == 1 {
		game.Player1Turn = true
	} else {
		game.Player1Turn = false
	}
	err = game.OpponentLastGuess.UnmarshalJSON(gameMap[9])
	if err != nil {
		return err
	}
	return nil

}

type games map[uint64]game

func (games *games) MarshalJSON() (b []byte, err error) {
	gamesMap := make(map[uint64][]byte)

	for i, game := range *games {
		gamesMap[i], err = game.MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
	}
	return json.Marshal(gamesMap)
}
func (games *games) UnmarshalJSON(b []byte) (err error) {
	gamesMap := make(map[uint64][]byte)
	err = json.Unmarshal(b, &gamesMap)
	if err != nil {
		return err
	}
	for i, value := range gamesMap {
		var tempGame game
		err = tempGame.UnmarshalJSON(value)
		if err != nil {
			return err
		}
		(*games)[i] = tempGame
	}
	return nil
}
func (games *games) getGamesFromState() {
	gamesBytes := state.ReadBytesByKey("games")
	err := games.UnmarshalJSON(gamesBytes)
	if err != nil {
		panic(err)
	}
}
