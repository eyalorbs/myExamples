package main

//todo: find a way to add money
//todo find a way to validate the ships on a board
//todo develop the user interface

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/service"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
	"log"
	"math"
)

var PUBLIC = sdk.Export(startGame, getContractBalance, getUserBalance, getOpponentStatus, guess, updateHit, quitGame, getMyHits, approveBoard, checkIfWon)
var SYSTEM = sdk.Export(_init, PanicIfSignerPlaying, PanicIfSignerNotPlaying, HandleMoney, ApproveShipsOnBoard)

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

func HandleMoney() {
	//this is a temporary function that does nothing.
	//it will be replaced once I know how to deal with tokens
}

//public functions

//the user needs to manually approve the function 10 tokens
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
		defaultPlayerHits, defaultPlayerGuesses, defaultPlayerGuesses, turnDefault, defaultLastGuess, 3, 3)

	//check which index is free
	for i := uint64(1); i < math.MaxUint64; i++ {
		if _, ok := games[i]; !ok {
			//add the newGame go the games
			games[i] = newGame
			//update state
			games.updateState()

			//update the state: player's address are read as the index of the game they are playing
			state.WriteUint64ByAddress(address.GetSignerAddress(), i)
			state.ClearByAddress(player2)
			state.WriteUint64ByAddress(player2, i)
			HandleMoney() //DON'T FORGET TO REPLACE THIS BY A REAL FUNCTION
			return
		}
	}
	panic("no more room for more games, this game is very very very successful")
}

func guess(x, y uint32) {
	PanicIfSignerNotPlaying()
	//create a coordinate
	coo := coordinate{uint8(x), uint8(y)}
	//make sure the coordinate is valid
	coo.validateGuessCoordinates()
	//get the relevant game
	signerAddress := address.GetSignerAddress()
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(signerAddress)]
	//validate that the player is only playing if it's his turn
	relevantGame.panicIfNotTurn()
	//update last game
	relevantGame.OpponentLastGuess = coo
	//update the guesses
	relevantGame.updateGuesses(coo)
	//update state:
	games[state.ReadUint64ByAddress(signerAddress)] = relevantGame
	games.updateState()

}
func (coo *coordinate) validateGuessCoordinates() {
	//make sure user is playing
	PanicIfSignerPlaying()
	if 9 < coo.X || 9 < coo.Y {
		panic("guess out of range")
	}
	//get the relevant game
	signerAddress := address.GetSignerAddress()
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(signerAddress)]
	//make sure the spot wasn't guessed
	if bytes.Equal(signerAddress, relevantGame.Player1) {
		if relevantGame.Player1Guesses.exists(*coo) {
			panic("you already guessed this spot")
		}
	}
	if bytes.Equal(signerAddress, relevantGame.Player2) {
		if relevantGame.Player2Guesses.exists(*coo) {
			panic("you already guessed this spot")
		}
	}

}

func getOpponentStatus() (x, y uint32) {
	signerAddress := address.GetSignerAddress()
	PanicIfSignerNotPlaying()
	//get games and index of player
	index := state.ReadUint64ByAddress(signerAddress)
	games := make(games)
	games.getGamesFromState()
	//if the opponent hasn't played yed, panic. otherwise return the opponent guess
	opponentGuess := games[index].OpponentLastGuess
	EmptyGuess := coordinate{}
	if opponentGuess != EmptyGuess {
		panic("opponent didn't play yet")
	}
	return uint32(opponentGuess.X), uint32(opponentGuess.Y)
}

func updateHit(hit uint32) {
	PanicIfSignerNotPlaying()
	signerAddress := address.GetSignerAddress()
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(signerAddress)]
	relevantGame.panicIfTurn()
	if hit == 0 {
		if relevantGame.Player1Turn {
			relevantGame.Player1Guesses.opponentResponses = append(relevantGame.Player1Guesses.opponentResponses, false)
			relevantGame.Player1Turn = !relevantGame.Player1Turn
		} else if !relevantGame.Player1Turn {
			relevantGame.Player2Guesses.opponentResponses = append(relevantGame.Player2Guesses.opponentResponses, false)
			relevantGame.Player1Turn = !relevantGame.Player1Turn
		}
	} else {
		if relevantGame.Player1Turn {
			relevantGame.Player1Guesses.opponentResponses = append(relevantGame.Player1Guesses.opponentResponses, true)
			relevantGame.Player1Turn = !relevantGame.Player1Turn
			relevantGame.Player1Hits += 1
		} else if !relevantGame.Player1Turn {
			relevantGame.Player2Guesses.opponentResponses = append(relevantGame.Player2Guesses.opponentResponses, true)
			relevantGame.Player1Turn = !relevantGame.Player1Turn
			relevantGame.Player2Hits += 1
		}
	}
}

func getMyHits() (hits uint32) {
	PanicIfSignerNotPlaying()
	signerAddress := address.GetSignerAddress()
	games := make(games)
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(signerAddress)]
	if bytes.Equal(signerAddress, relevantGame.Player1) {
		return uint32(relevantGame.Player1Hits)
	}
	if bytes.Equal(signerAddress, relevantGame.Player2) {
		return uint32(relevantGame.Player2Hits)
	}
	panic("you are not registered for this game")

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

func quitGame() {
	PanicIfSignerNotPlaying()

	HandleMoney() //DON'T FORGET TO REPLACE THIS BY A REAL FUNCTION
	signerAddress := address.GetSignerAddress()
	games := make(games)
	games.getGamesFromState()
	game := games[state.ReadUint64ByAddress(signerAddress)]

	state.ClearByAddress(game.Player1)
	state.ClearByAddress(game.Player2)
	delete(games, state.ReadUint64ByAddress(signerAddress))
	games.updateState()

}

func approveBoard(secretKey string, shaBoard []byte, board []byte) {
	PanicIfSignerNotPlaying()

	//get the relevant game
	signerAddress := address.GetSignerAddress()
	games := make(games)
	games.getGamesFromState()
	game := games[state.ReadUint64ByAddress(signerAddress)]

	//get the board the player claims to have
	secretBoard := secretBoard{}
	err := secretBoard.UnmarshalJSON(board)
	if err != nil {
		panic(err)
	}
	//calculate the sha with the secret key
	realSha := secretBoard.sha256(secretKey)

	//if the player's sha and the calculated sha are different, do not approve board
	if !bytes.Equal(realSha, shaBoard) {
		if bytes.Equal(game.Player1, signerAddress) {
			game.board1Approved = 2
			games[state.ReadUint64ByAddress(signerAddress)] = game
			games.updateState()
		} else if bytes.Equal(game.Player2, signerAddress) {
			game.board2Approved = 2
			games[state.ReadUint64ByAddress(signerAddress)] = game
			games.updateState()
		} else {
			panic("you are not registered for this game")
		}
	}

	if bytes.Equal(game.Player1, signerAddress) {
		//if there ins't the same amount of guesses as responses the game cannot be over
		if len(game.Player2Guesses.opponentResponses) != len(game.Player2Guesses.playerGuesses) {
			panic("the game cannot be over, there isn't the same amount of guesses as responses")
		}
		//compare between the gusses and answer
		for i := 0; i < len(game.Player2Guesses.playerGuesses); i++ {
			for j := 0; j < len(secretBoard); j++ {
				if secretBoard[j].coo == game.Player2Guesses.playerGuesses[i] {
					if secretBoard[j].ship != game.Player2Guesses.opponentResponses[i] {
						game.board1Approved = 2
						games[state.ReadUint64ByAddress(signerAddress)] = game
						games.updateState()
					}
				}
			}
		}
	} else if bytes.Equal(game.Player2, signerAddress) {
		if len(game.Player1Guesses.opponentResponses) != len(game.Player1Guesses.playerGuesses) {
			panic("the game cannot be over, there isn't the same amount of guesses as responses")
		}
		//compare between the gusses and answer
		for i := 0; i < len(game.Player1Guesses.playerGuesses); i++ {
			for j := 0; j < len(secretBoard); j++ {
				if secretBoard[j].coo == game.Player1Guesses.playerGuesses[i] {
					if secretBoard[j].ship != game.Player1Guesses.opponentResponses[i] {
						game.board2Approved = 2
						games[state.ReadUint64ByAddress(signerAddress)] = game
						games.updateState()
					}
				}
			}
		}
	} else {
		panic("you are not registered for this game")
	}

	if bytes.Equal(game.Player1, signerAddress) {
		game.board1Approved = ApproveShipsOnBoard(secretBoard)
		games[state.ReadUint64ByAddress(signerAddress)] = game
		games.updateState()
	} else if bytes.Equal(game.Player2, signerAddress) {
		game.board2Approved = ApproveShipsOnBoard(secretBoard)
		games[state.ReadUint64ByAddress(signerAddress)] = game
		games.updateState()
	} else {
		panic("you are not registered for this game")
	}

}
func ApproveShipsOnBoard(secretBoard secretBoard) (approved uint8) {
	var board [10][10]bool
	shipCoordinates := []coordinate{}
	for i := 0; i < len(secretBoard); i++ {
		board[i/10][i%10] = secretBoard[i].ship
		if secretBoard[i].ship {
			shipCoordinates = append(shipCoordinates, coordinate{uint8(i / 10), uint8(i % 10)})
		}
	}
	if len(shipCoordinates) != 17 {
		return 2
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////////////find out how to prove that a board is ok
	////////////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////////////
	return 1
}
func checkIfWon() (player1, player2 bool) {
	PanicIfSignerNotPlaying()
	//get games and relevant game
	games := make(games)
	signerAddress := address.GetSignerAddress()
	games.getGamesFromState()
	relevantGame := games[state.ReadUint64ByAddress(signerAddress)]
	//if neither player's have enough points, do not proceed to check who won
	if relevantGame.Player1Hits != 17 && relevantGame.Player2Hits != 17 {
		panic("neither player have enough points to win")
	}
	//if one or more of the players did not approve their board do not proceed
	if relevantGame.board1Approved == 3 || relevantGame.board2Approved == 3 {
		panic("both players need to prove their board is ok")
	}

	if relevantGame.board1Approved == 2 {
		if relevantGame.board2Approved == 2 {
			//if both player's cheated, none of them get their money
			return false, false
		} else {
			//if only player1 cheated, player2 wins
			return false, true
		}
	} else {
		//if only player2 cheated, player1 wins
		if relevantGame.board2Approved == 2 {
			return true, false
		}
	}

	//if neither player's cheated, whoever has 17 hits wins
	if relevantGame.Player1Hits == 17 {
		return true, false
	} else {
		return false, true
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
	signerAddress := address.GetSignerAddress()
	if bytes.Equal(signerAddress, game.Player1) {
		if !game.Player1Turn {
			panic("not your turn, wait for your turn")
		}
	} else if bytes.Equal(signerAddress, game.Player2) {
		if game.Player1Turn {
			panic("not your turn, wait for your turn")
		}
	} else {
		panic("you are not registered for this game")
	}
}
func (game *game) panicIfTurn() {
	signerAddress := address.GetSignerAddress()
	if bytes.Equal(signerAddress, game.Player1) {
		if game.Player1Turn {
			panic("it is your turn, you cannot validate ship")
		}
	} else if bytes.Equal(signerAddress, game.Player2) {
		if !game.Player1Turn {
			panic("it is your turn, you cannot validate ship")
		}
	} else {
		panic("you are not registered for this game")
	}
}

func (game *game) updateGuesses(coo coordinate) {
	signerAddress := address.GetSignerAddress()
	if bytes.Equal(signerAddress, game.Player1) {
		game.Player1Guesses.playerGuesses = append(game.Player1Guesses.playerGuesses, coo)
	} else if bytes.Equal(signerAddress, game.Player2) {
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

type box struct {
	coo  coordinate
	ship bool
}

func (box *box) new(coo coordinate, ship bool) {
	box.coo = coo
	box.ship = ship
}
func (box *box) MarshalJSON() (b []byte, err error) {
	boxMap := make(map[uint8][]byte)
	b, err = box.coo.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}
	boxMap[0] = b
	if box.ship {
		boxMap[1] = []byte{1}
	} else {
		boxMap[1] = []byte{0}
	}
	return json.Marshal(boxMap)
}
func (box *box) UnmarshalJSON(b []byte) (err error) {
	boxMap := make(map[uint8][]byte)
	err = json.Unmarshal(b, &boxMap)
	if err != nil {
		return err
	}
	coo := coordinate{}
	err = coo.UnmarshalJSON(boxMap[0])
	if err != nil {
		return err
	}
	box.coo = coo
	box.ship = bytes.Equal(boxMap[1], []byte{1})
	return nil
}

type secretBoard []box

func (boxes *secretBoard) MarshalJSON() (b []byte, err error) {
	secretBoardMap := make(map[uint8][]byte)
	for i, val := range *boxes {
		b, err = val.MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
		secretBoardMap[uint8(i)] = b
	}
	return json.Marshal(secretBoardMap)
}
func (boxes *secretBoard) UnmarshalJSON(b []byte) (err error) {
	secretBoardMap := make(map[uint8][]byte)
	err = json.Unmarshal(b, &secretBoardMap)
	if err != nil {
		return err
	}
	var box box
	for i := 0; i < len(secretBoardMap); i++ {
		err = box.UnmarshalJSON(secretBoardMap[uint8(i)])
		if err != nil {
			return err
		}
		*boxes = append(*boxes, box)
	}
	return nil
}
func (boxes *secretBoard) sha256(sk string) (sha []byte) {
	h := hmac.New(sha256.New, []byte(sk))
	b, err := boxes.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	h.Write(b)
	return h.Sum(nil)

}
