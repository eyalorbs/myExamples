package main

//todo: find a way to add money
//todo develop the user interface

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/events"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/service"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
	"math"
)

var PUBLIC = sdk.Export(startGame, getOpponentStatus, guess, updateHit, quitGame, getMyHits, approveBoard, checkIfWon, checkIfInGame, didOpponentUpdateHit, getIfWonLastGame)
var SYSTEM = sdk.Export(_init, panicIfCallerPlaying, PanicIfCallerNotPlaying, shipsOk, didLie, InGame, PanicIfOneOfPlayersApprovedBoard)

//helper coin contract name: token-bridge

func _init() {
	//create an empty instance of games
	games := make(games)
	b, err := games.MarshalJSON()
	if err != nil {
		panic(err)
	}
	//update games and the waiting pool to the state
	state.WriteBytesByKey("games", b)
	state.WriteBytesByKey("waitingPool", []byte{})
}

//helper function: panics if one of the players approved the board
func PanicIfOneOfPlayersApprovedBoard() {
	//get the games from the state
	games := make(games)
	games.getGamesFromState()
	//get the relevant game to the caller
	relevantGame := games[state.ReadUint64ByAddress(address.GetCallerAddress())]
	//check if one of the players approved their board
	if relevantGame.board1Approved != 3 || relevantGame.board2Approved != 3 {
		panic("both players need to approve their board")
	}
}

//checks if the caller has won the last game
//it calls a contract that keeps track
//1 means that there is no previous game
//2 means that the player lost
//3 means that the player won
func getIfWonLastGame() (didWin uint32) {
	panicIfCallerPlaying()
	return service.CallMethod("winnerContract", "getWinner", address.GetCallerAddress())[0].(uint32)
}

//helper function: panics if the  caller is playing
func panicIfCallerPlaying() {
	callerAddress := address.GetCallerAddress()

	//if there are 256 bits(a hashedBoard) than the player is playing
	if len(state.ReadBytesByAddress(callerAddress)) == 32 {
		panic("player is already in pool")
	}

	//if the index isn't 0, the player is playing
	index := state.ReadUint64ByAddress(callerAddress)
	if index != 0 {
		panic("player is already playing")
	}

}

//helper function: panics if the caller isn't playing
func PanicIfCallerNotPlaying() {
	callerAddress := address.GetCallerAddress()

	//if there are 256 bits(a hashedBoard) than the player is playing
	if len(state.ReadBytesByAddress(callerAddress)) == 32 {
		panic("player is in pool")
	}

	//if the index isn't 0, the player is playing
	index := state.ReadUint64ByAddress(callerAddress)
	if index == 0 {
		panic("player is not playing a game")
	}
}

//public functions

//event for startGame
func startGameEvent(event string) {

}

func startGame(hashedBoard []byte) {
	//makes sure that the size of the hashed board is valid
	if len(hashedBoard) != 32 {
		panic("hashed board must be 256 bits")
	}
	//transfers money to contract with helper contract
	service.CallMethod("tokenBridge", "transferToContract", address.GetCallerAddress(), uint64(5))
	//get the games from the state
	games := make(games)
	games.getGamesFromState()
	//make sure the player isn't playing
	panicIfCallerPlaying()

	//if there isn't anyone in the pool, add the player to the pool and save his hashed ships
	if bytes.Equal(state.ReadBytesByKey("waitingPool"), []byte{}) {
		state.WriteBytesByKey("waitingPool", address.GetCallerAddress())
		state.WriteBytesByAddress(address.GetCallerAddress(), hashedBoard)

		//send an event
		events.EmitEvent(startGameEvent, "added to pool")
		return
	}
	//if there is someone in the waiting pool
	//player1 is the caller address
	player1 := address.GetCallerAddress()

	//get player2Address from the pool
	player2 := state.ReadBytesByKey("waitingPool")[:20]
	//update the waiting pool
	state.WriteBytesByKey("waitingPool", state.ReadBytesByKey("waitingPool")[20:])

	//get the values necessary to start a new game:
	board1Hashed := hashedBoard
	//player2 hashed board is read from the state
	board2Hashed := state.ReadBytesByAddress(player2)
	defaultPlayerHits := uint8(0)
	defaultPlayerGuesses := guesses{}
	turnDefault := true
	defaultLastGuess := coordinate{}
	var newGame game
	newGame.new(player1, player2, board1Hashed, board2Hashed, defaultPlayerHits,
		defaultPlayerHits, defaultPlayerGuesses, defaultPlayerGuesses, turnDefault, defaultLastGuess, 3, 3)

	//check which key is free and add the new game to the free key
	for i := uint64(1); i < math.MaxUint64; i++ {
		if _, ok := games[i]; !ok {
			//add the newGame go the games
			games[i] = newGame
			//update state
			games.updateState()

			//update the state: player's address are read as the index of the game they are playing
			state.WriteUint64ByAddress(player1, i)
			state.ClearByAddress(player2)
			state.WriteUint64ByAddress(player2, i)
			events.EmitEvent(startGameEvent, "started new game")
			return
		}
	}
	panic("no more room for more games, this game is very very popular")
}

//event for checkInGame function
func InGame(inGame string, x, y uint32) {}
func checkIfInGame() {
	//get the caller address
	callerAddress := address.GetCallerAddress()
	//get the games from the state
	games := make(games)
	games.getGamesFromState()
	//get the relevant game for the caller
	relevantGame := games[state.ReadUint64ByAddress(callerAddress)]
	//if there are 256 bits(a hashedBoard) than the player is playing
	if len(state.ReadBytesByAddress(callerAddress)) == 32 {
		events.EmitEvent(InGame, "player still in pool", uint32(10), uint32(10))
		return
	}

	//if the index isn't 0, the player is playing
	index := state.ReadUint64ByAddress(callerAddress)
	if index != 0 {
		events.EmitEvent(InGame, "player is in a game", uint32(relevantGame.OpponentLastGuess.X), uint32(relevantGame.OpponentLastGuess.Y))
		return
	}

	events.EmitEvent(InGame, "player is neither in pool nor in a game", uint32(10), uint32(10))
}

//event for guess function
func guessEvent(response string) {}
func guess(x, y uint32) {
	//panic if the caller is not playing
	PanicIfCallerNotPlaying()
	//panic if the one of the player's approved their board.
	PanicIfOneOfPlayersApprovedBoard()

	//create a coordinate
	coo := coordinate{uint8(x), uint8(y)}
	//make sure the coordinate is valid
	coo.validateGuessCoordinates()
	//get the relevant game
	callerAddress := address.GetCallerAddress()
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(callerAddress)]
	//validate that the player is only playing if it's his turn
	relevantGame.panicIfNotTurn()

	if bytes.Equal(relevantGame.Player1, callerAddress) {
		//if the player needs to approve the board because the game has ended tell him:
		if relevantGame.board2Approved != 3 {
			events.EmitEvent(guessEvent, "you need to approve your board")
		} else if relevantGame.board1Approved != 3 {
			events.EmitEvent(guessEvent, "you already approved your board, we are in the endgame now")
		}
		relevantGame.Player1Guesses.playerGuesses = append(relevantGame.Player1Guesses.playerGuesses, coo)

	} else if bytes.Equal(relevantGame.Player2, callerAddress) {
		if relevantGame.board1Approved != 3 {
			events.EmitEvent(guessEvent, "you need to approve your board")

		} else if relevantGame.board2Approved != 3 {
			events.EmitEvent(guessEvent, "you already approved your board, we are in the endgame now")
		}
		relevantGame.Player2Guesses.playerGuesses = append(relevantGame.Player2Guesses.playerGuesses, coo)
	} else {
		panic("you are not registered for this game")
	}
	//update last game
	relevantGame.OpponentLastGuess = coo
	//update state:

	games[state.ReadUint64ByAddress(callerAddress)] = relevantGame
	games.updateState()
	events.EmitEvent(guessEvent, "guess submitted")

	// check if the player needs to approve board
	if relevantGame.Player1Hits == 17 || relevantGame.Player2Hits == 17 {
		events.EmitEvent(guessEvent, "approve your board")
	}
}
func (coo *coordinate) validateGuessCoordinates() {
	//make sure user is playing
	PanicIfCallerNotPlaying()
	if coo.X == 0 || 10 < coo.X || 10 < coo.Y || coo.Y == 0 {
		panic("guess out of range")
	}
	//get the relevant game
	callerAddress := address.GetCallerAddress()
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(callerAddress)]
	//make sure the spot wasn't guessed
	if bytes.Equal(callerAddress, relevantGame.Player1) {
		if relevantGame.Player1Guesses.exists(*coo) {
			panic("you already guessed this spot")
		}
	}
	if bytes.Equal(callerAddress, relevantGame.Player2) {
		if relevantGame.Player2Guesses.exists(*coo) {
			panic("you already guessed this spot")
		}
	}

}

//returns the last guess of the opponent
func getOpponentStatus() (x, y uint32) {
	//get the caller's address
	callerAddress := address.GetCallerAddress()
	//panic if the caller is not playing
	PanicIfCallerNotPlaying()
	//panic if one of the player's approved their board
	PanicIfOneOfPlayersApprovedBoard()
	//get the relevant game
	index := state.ReadUint64ByAddress(callerAddress)
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[index]
	//panic if it is the player's turn
	if bytes.Equal(relevantGame.Player1, callerAddress) {
		if relevantGame.Player1Turn {
			panic("it is your turn, you are not the one who is supposed to validate it")
		}
	} else if bytes.Equal(relevantGame.Player2, callerAddress) {
		if !relevantGame.Player1Turn {
			panic("it is your turn, you are not the one who is supposed to validate it")
		}
	}

	//if the opponent hasn't played yet, panic. otherwise return the opponent guess
	opponentGuess := relevantGame.OpponentLastGuess
	EmptyGuess := coordinate{}
	if opponentGuess == EmptyGuess {
		panic("opponent didn't play yet")
	}
	return uint32(opponentGuess.X), uint32(opponentGuess.Y)
}

func updateHitEmit(feedback string) {}
func updateHit(hit uint32) {
	PanicIfCallerNotPlaying()
	//get the relevant game
	callerAddress := address.GetCallerAddress()
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(callerAddress)]

	relevantGame.panicIfTurn()

	//only let player update if he needs to
	if hit == 0 {
		if relevantGame.Player1Turn {
			if len(relevantGame.Player1Guesses.opponentResponses) == len(relevantGame.Player1Guesses.playerGuesses) {
				panic("there is no need to update")
			}
			//add the response to the game and change the turn
			relevantGame.Player1Guesses.opponentResponses = append(relevantGame.Player1Guesses.opponentResponses, false)
			relevantGame.Player1Turn = !relevantGame.Player1Turn
		} else {
			if len(relevantGame.Player2Guesses.opponentResponses) == len(relevantGame.Player2Guesses.playerGuesses) {
				panic("there is no need to update")
			}
			//add the response to the game and change the turn
			relevantGame.Player2Guesses.opponentResponses = append(relevantGame.Player2Guesses.opponentResponses, false)
			relevantGame.Player1Turn = !relevantGame.Player1Turn
		}
	} else {
		if relevantGame.Player1Turn {
			if len(relevantGame.Player1Guesses.opponentResponses) == len(relevantGame.Player1Guesses.playerGuesses) {
				panic("there is no need to update")
			}
			//add 1 to the player's hits, add the response to the game and change the turn
			relevantGame.Player1Hits += 1
			relevantGame.Player1Guesses.opponentResponses = append(relevantGame.Player1Guesses.opponentResponses, true)
			relevantGame.Player1Turn = !relevantGame.Player1Turn

		} else {
			if len(relevantGame.Player2Guesses.opponentResponses) == len(relevantGame.Player2Guesses.playerGuesses) {
				panic("there is no need to update")
			}
			//add 1 to the player's hits, add the response to the game and change the turn
			relevantGame.Player2Hits += 1
			relevantGame.Player2Guesses.opponentResponses = append(relevantGame.Player2Guesses.opponentResponses, true)
			relevantGame.Player1Turn = !relevantGame.Player1Turn

		}
	}
	relevantGame.OpponentLastGuess = coordinate{}
	//update state and games
	games[state.ReadUint64ByAddress(callerAddress)] = relevantGame
	games.updateState()
}

//event for didOpponentUpdateHit
func hitUpdatedEvent(hit uint32) {}

//tells the player if his opponent updated his guess
func didOpponentUpdateHit() {
	PanicIfCallerNotPlaying()
	callerAddress := address.GetCallerAddress()
	//get the relevant game from the state
	games := make(games)
	games.getGamesFromState()
	index := state.ReadUint64ByAddress(callerAddress)
	relevantGame := games[index]
	//if the length of the guesses slice and the responses slice is equal the hit has been updated
	if bytes.Equal(relevantGame.Player1, callerAddress) {
		if len(relevantGame.Player1Guesses.playerGuesses) == len(relevantGame.Player1Guesses.opponentResponses) {
			events.EmitEvent(hitUpdatedEvent, uint32(1))
			return
		}
		events.EmitEvent(hitUpdatedEvent, uint32(0))
	} else if bytes.Equal(relevantGame.Player2, callerAddress) {
		if len(relevantGame.Player1Guesses.playerGuesses) == len(relevantGame.Player1Guesses.opponentResponses) {
			events.EmitEvent(hitUpdatedEvent, uint32(1))
			return
		}
		events.EmitEvent(hitUpdatedEvent, uint32(0))
	} else {
		panic("it is not your turn")
	}

}

//returns the number of hits
func getMyHits() (hits uint32) {
	PanicIfCallerNotPlaying()
	callerAddress := address.GetCallerAddress()
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(callerAddress)]
	if bytes.Equal(callerAddress, relevantGame.Player1) {
		return uint32(relevantGame.Player1Hits)
	}
	if bytes.Equal(callerAddress, relevantGame.Player2) {
		return uint32(relevantGame.Player2Hits)
	}
	panic("you are not registered for this game")

}

//quits game
func quitGame() {
	PanicIfCallerNotPlaying()
	callerAddress := address.GetCallerAddress()
	//get the relevant game from the state
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(callerAddress)]
	//give the winnings to the player who didn't quit
	bridgeAddress, _ := hex.DecodeString("8fef8b50287ce37cd9a738393822c20ba0b7cf2d")
	service.CallMethod("token", "approve", bridgeAddress, uint64(9))

	if bytes.Equal(callerAddress, relevantGame.Player1) {
		service.CallMethod("tokenBridge", "transferToWinner", relevantGame.Player2, uint64(9))
	} else {
		service.CallMethod("tokenBridge", "transferToWinner", relevantGame.Player1, uint64(9))
	}
	//update the games and the state
	delete(games, state.ReadUint64ByAddress(callerAddress))
	state.ClearByAddress(relevantGame.Player1)
	state.ClearByAddress(relevantGame.Player2)
	games.updateState()

}

//returns if the ships locations are ok and their coordinates
func shipsOk(boats ships) (ok bool, shipCoordinates []coordinate) {
	//if there aren't 5 boars return false
	if len(boats) != 5 {
		return false, nil
	}
	//go over each boat
	for _, val := range boats {
		//check if boat is diagonal
		if val.headCoordinates.X != val.tailCoordinates.X && val.headCoordinates.Y != val.tailCoordinates.Y {
			return false, nil
		}
		//check if coordinates are in range
		if 10 < val.headCoordinates.X || 10 < val.headCoordinates.Y || 10 < val.tailCoordinates.X || 10 < val.tailCoordinates.Y {
			return false, nil
		}

		//check if length is ok
		var length uint8
		switch val.name {
		case "Carrier":
			length = 5

		case "Battleship":
			length = 4

		case "Cruiser":
			length = 3

		case "Submarine":
			length = 3

		case "Destroyer":
			length = 2

		default:
			return false, nil
		}
		if uint8(math.Abs(float64(val.headCoordinates.X)-float64(val.tailCoordinates.X))) != length && uint8(math.Abs(float64(val.headCoordinates.Y)-float64(val.tailCoordinates.Y))) != length {
			return false, nil
		}

		//add all of the coordinates to a slice and return it, if there is overlap return false
		for i := uint8(math.Min(float64(val.headCoordinates.X), float64(val.tailCoordinates.X))); i <= uint8(math.Max(float64(val.headCoordinates.X), float64(val.tailCoordinates.X))); i++ {
			for j := uint8(math.Min(float64(val.headCoordinates.Y), float64(val.tailCoordinates.Y))); j <= uint8(math.Max(float64(val.headCoordinates.Y), float64(val.tailCoordinates.Y))); j++ {
				for _, coor := range shipCoordinates {
					currentCoo := coordinate{i, j}
					if currentCoo == coor {

						return false, nil
					}
					shipCoordinates = append(shipCoordinates, currentCoo)
				}
			}
		}
	}
	return true, shipCoordinates
}

//check if the player lied about the opponent's hits
func didLie(shipCoordinates []coordinate, relevantGame game) (lied bool) {
	callerAddress := address.GetCallerAddress()
	//if player 1
	if bytes.Equal(callerAddress, relevantGame.Player1) {
		//go over the guesses
		for i, val := range relevantGame.Player1Guesses.playerGuesses {
			found := false
			//if they are in the coordinates, check if the opponent said that they were hit
			if relevantGame.Player1Guesses.opponentResponses[i] {
				for _, coo := range shipCoordinates {
					if val == coo {
						found = true
						break
					}
				}
				if !found {
					return false
				}

			}
		}
		//same for player 2
	} else if bytes.Equal(callerAddress, relevantGame.Player2) {
		for i, val := range relevantGame.Player2Guesses.playerGuesses {
			found := false
			if relevantGame.Player2Guesses.opponentResponses[i] {
				for _, coo := range shipCoordinates {
					if val == coo {
						found = true
						break
					}
				}
				if !found {
					return false
				}

			}
		}

	} else {
		panic("you are not registered for this game")
	}

	return true

}

func approveBoard(secretKey string, Marshaledships []byte) {
	PanicIfCallerNotPlaying()

	//get the relevant game
	callerAddress := address.GetCallerAddress()
	games := make(games)
	games.getGamesFromState()
	game := games[state.ReadUint64ByAddress(callerAddress)]

	approve := true
	//get the board the player claims to have
	boats := ships{}
	err := boats.UnmarshalJSON(Marshaledships)
	if err != nil {
		approve = false
	}
	//calculate the sha with the secret key
	realSha, err := boats.sha256(secretKey)
	if err != nil {
		approve = false
	}
	ok, coo := shipsOk(boats)

	if bytes.Equal(game.Player1, callerAddress) {
		shaShips := game.Board1Hashed
		if !approve {
			game.board1Approved = 2
		} else if !bytes.Equal(realSha, shaShips) {
			game.board1Approved = 2
		} else if !ok {
			game.board1Approved = 2
		} else if didLie(coo, game) {
			game.board1Approved = 2
		}
		game.board1Approved = 1
		games[state.ReadUint64ByAddress(callerAddress)] = game
		games.updateState()

	} else if bytes.Equal(game.Player2, callerAddress) {
		shaShips := game.Board2Hashed
		if !approve {
			game.board2Approved = 2
		} else if !bytes.Equal(realSha, shaShips) {
			game.board2Approved = 2
		} else if !ok {
			game.board2Approved = 2
		} else if didLie(coo, game) {
			game.board2Approved = 2
		}
		game.board2Approved = 1
		games[state.ReadUint64ByAddress(callerAddress)] = game
		games.updateState()

	} else {
		panic("you are not registered for this game")
	}

}

func checkIfWon() {
	PanicIfCallerNotPlaying()
	//get games and relevant game
	games := make(games)
	callerAddress := address.GetCallerAddress()
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(callerAddress)]
	//if neither player's have enough points, do not proceed to check who won
	if relevantGame.Player1Hits != 17 && relevantGame.Player2Hits != 17 {
		panic("neither player have enough points to win")
	}
	//if one or more of the players did not approve their board do not proceed
	if relevantGame.board1Approved == 3 || relevantGame.board2Approved == 3 {
		panic("both players need to prove their board is ok")
	}
	player1Address := relevantGame.Player1
	player2Address := relevantGame.Player2
	board1Approved := relevantGame.board1Approved
	board2Approved := relevantGame.board2Approved
	player1Hits := relevantGame.Player1Hits

	delete(games, state.ReadUint64ByAddress(callerAddress))
	state.ClearByAddress(relevantGame.Player1)
	state.ClearByAddress(relevantGame.Player2)
	games.updateState()
	bridgeAddress, _ := hex.DecodeString("8fef8b50287ce37cd9a738393822c20ba0b7cf2d")

	if board1Approved == 2 {
		if board2Approved == 2 {
			//if both player's cheated, none of them get their money
			service.CallMethod("winnerContract", "writeWinner", player1Address, uint32(1))
			service.CallMethod("winnerContract", "writeWinner", player2Address, uint32(1))
		} else {
			//if only player1 cheated, player2 wins
			service.CallMethod("winnerContract", "writeWinner", player1Address, uint32(1))
			service.CallMethod("winnerContract", "writeWinner", player2Address, uint32(2))
			//give the tokens to the winner
			service.CallMethod("token", "approve", bridgeAddress, uint64(9))
			service.CallMethod("tokenBridge", "transferToWinner", player2Address, uint64(9))
		}
	} else {
		//if only player2 cheated, player1 wins
		if board2Approved == 2 {
			service.CallMethod("winnerContract", "writeWinner", player1Address, uint32(2))
			service.CallMethod("winnerContract", "writeWinner", player2Address, uint32(1))
			//give the tokens to the winner
			service.CallMethod("token", "approve", bridgeAddress, uint64(9))
			service.CallMethod("tokenBridge", "transferToWinner", player1Address, uint64(9))
		}
	}

	//if neither player's cheated, whoever has 17 hits wins
	if player1Hits == 17 {
		service.CallMethod("winnerContract", "writeWinner", player1Address, uint32(2))
		service.CallMethod("winnerContract", "writeWinner", player2Address, uint32(1))
		//give the tokens to the winner
		service.CallMethod("token", "approve", bridgeAddress, uint64(9))
		service.CallMethod("tokenBridge", "transferToWinner", player1Address, uint64(9))
	} else {
		service.CallMethod("winnerContract", "writeWinner", player1Address, uint32(1))
		service.CallMethod("winnerContract", "writeWinner", player2Address, uint32(2))
		//give the tokens to the winner
		service.CallMethod("token", "approve", bridgeAddress, uint64(9))
		service.CallMethod("tokenBridge", "transferToWinner", player2Address, uint64(9))
	}

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

type guesses struct {
	playerGuesses     []coordinate
	opponentResponses []bool
}

func (guesses *guesses) MarshalJSON() (b []byte, err error) {
	lengthGuesses := uint8(len(guesses.playerGuesses))
	guessesMap := make(map[uint8][]byte)
	for i := uint8(0); i < lengthGuesses; i++ {
		b, err := guesses.playerGuesses[i].MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
		guessesMap[i] = b
	}
	lengthResponses := uint8(len(guesses.opponentResponses))
	for i := uint8(0); i < lengthResponses; i++ {
		if guesses.opponentResponses[i] {
			guessesMap[i] = append(guessesMap[i], byte(1))
		} else {
			guessesMap[i] = append(guessesMap[i], byte(0))
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
	for i := 0; i < len(guessesMap); i++ {
		var guess coordinate
		val := guessesMap[uint8(i)]
		err = guess.UnmarshalJSON(val)
		if err != nil {
			err = guess.UnmarshalJSON(val[:len(val)-1])
			if err != nil {
				return err
			}
			guesses.playerGuesses = append(guesses.playerGuesses, guess)
			guesses.opponentResponses = append(guesses.opponentResponses, val[len(val)-1] == 1)
		} else {
			guesses.playerGuesses = append(guesses.playerGuesses, guess)
		}
	}
	return nil
}
func (guesses *guesses) exists(coo coordinate) (exists bool) {
	//return true if the coordinate exists in the previous guesses
	for _, value := range guesses.playerGuesses {
		if value == coo {
			return true
		}
	}
	return false
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

	//1 means approved, 2 means not approved, 3 means didn't check
	board1Approved uint8
	board2Approved uint8
}

func (game *game) new(player1, player2, board1Hashed, board2Hashed []byte, player1Hits, player2Hits uint8,
	player1Guesses, player2Guesses guesses, player1Turn bool, opponentLastGuess coordinate, board1Approved, board2Approved uint8) {
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
	game.board1Approved = board1Approved
	game.board2Approved = board2Approved
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
	gameMap[10] = []byte{game.board1Approved}
	gameMap[11] = []byte{game.board2Approved}
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
	game.board1Approved = gameMap[10][0]
	game.board2Approved = gameMap[11][0]

	return nil

}
func (game *game) panicIfNotTurn() {
	callerAddress := address.GetCallerAddress()
	if bytes.Equal(callerAddress, game.Player1) {
		if !game.Player1Turn {
			panic("not your turn, wait for your turn")
		}
	} else if bytes.Equal(callerAddress, game.Player2) {
		if game.Player1Turn {
			panic("not your turn, wait for your turn")
		}
	} else {
		panic("you are not registered for this game")
	}
}
func (game *game) panicIfTurn() {
	callerAddress := address.GetCallerAddress()
	if bytes.Equal(callerAddress, game.Player1) {
		if game.Player1Turn {
			panic("it is your turn, you cannot validate ship")
		}
	} else if bytes.Equal(callerAddress, game.Player2) {
		if !game.Player1Turn {
			panic("it is your turn, you cannot validate ship")
		}
	} else {
		panic("you are not registered for this game")
	}
}

func (game *game) updateGuesses(coo coordinate) {
	callerAddress := address.GetCallerAddress()
	if bytes.Equal(callerAddress, game.Player1) {
		game.Player1Guesses.playerGuesses = append(game.Player1Guesses.playerGuesses, coo)
	} else if bytes.Equal(callerAddress, game.Player2) {
		game.Player2Guesses.playerGuesses = append(game.Player2Guesses.playerGuesses, coo)
	} else {
		panic("you are not registered in this game")
	}

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
func (games *games) updateState() {
	b, err := games.MarshalJSON()
	if err != nil {
		panic(err)
	}
	state.WriteBytesByKey("games", b)
}

type ship struct {
	name            string
	headCoordinates coordinate
	tailCoordinates coordinate
}

func (boat *ship) new(name string, headX, headY, tailX, tailY uint8) {
	boat.name = name
	boat.headCoordinates = coordinate{headX, headY}
	boat.tailCoordinates = coordinate{tailX, tailY}
}
func (boat *ship) MarshalJSON() (b []byte, err error) {
	boatMap := make(map[uint8][]byte)
	boatMap[0] = []byte(boat.name)
	b, err = boat.headCoordinates.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}
	boatMap[1] = b

	b, err = boat.tailCoordinates.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}
	boatMap[2] = b
	return json.Marshal(boatMap)
}
func (boat *ship) UnmarshalJSON(b []byte) (err error) {
	boatMap := make(map[uint8][]byte)
	err = json.Unmarshal(b, &boatMap)
	if err != nil {
		return err
	}
	boat.name = string(boatMap[0])
	err = boat.headCoordinates.UnmarshalJSON(boatMap[1])
	if err != nil {
		return err
	}
	err = boat.tailCoordinates.UnmarshalJSON(boatMap[2])
	if err != nil {
		return err
	}
	return nil
}

type ships []ship

func (boats *ships) MarshalJSON() (b []byte, err error) {
	boatsMap := make(map[uint8][]byte)
	for i, val := range *boats {
		boatsMap[uint8(i)], err = val.MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
	}
	return json.Marshal(boatsMap)
}
func (boats *ships) UnmarshalJSON(b []byte) (err error) {
	boatsMap := make(map[uint8][]byte)
	err = json.Unmarshal(b, &boatsMap)
	if err != nil {
		return err
	}
	var temp ship
	for i := 0; i < len(boatsMap); i++ {
		err = temp.UnmarshalJSON(boatsMap[uint8(i)])
		if err != nil {
			return err
		}
		*boats = append(*boats, temp)
	}
	return nil
}
func (boats *ships) sha256(sk string) (sha []byte, err error) {
	h := hmac.New(sha256.New, []byte(sk))
	b, err := boats.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}
	h.Write(b)
	return h.Sum(nil), nil

}
