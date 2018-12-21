package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
)

var PUBLIC = sdk.Export(startGame, play, getBoard, checkIfWon, finishGame, quitPool)
var SYSTEM = sdk.Export(_init, uint64ToBytes, bytesToUint64)

//system functions
func _init() {
	//the addresses of the people waiting for a game
	state.WriteBytesByKey("waitingPool", []byte{})

	//the serialized games
	state.WriteBytesByKey("games", []byte{})

	//the free indexes (uint64 represented in a byte slice
	state.WriteBytesByKey("freeIndexes", []byte{})
}

func uint64ToBytes(num uint64) (b []byte) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	return b
}

func bytesToUint64(b []byte) (num uint64) {

	return binary.LittleEndian.Uint64(b)

}

//public functions and there helper methods
func startGame() (feedback string) {
	callerAddress := address.GetSignerAddress()

	//get the games from the state
	var games Games
	byteGames := state.ReadBytesByKey("games")
	//update the games
	err := games.UnmarshalJSON(byteGames)
	if err != nil {
		panic(err)
	}
	//make sure no one is starting a new game if he is already in a game
	if state.ReadUint64ByAddress(callerAddress) == 0 {
		if len(games) != 0 {
			if bytes.Equal(games[0].PlayerX, callerAddress) || bytes.Equal(games[0].PlayerO, callerAddress) {
				panic("you are already playing a game")
			}
		}
	}

	//if no one else is waiting, add the player to the waiting pool
	if len(state.ReadBytesByKey("waitingPool")) == 0 {
		state.WriteBytesByKey("waitingPool", append(state.ReadBytesByKey("waitingPool"), callerAddress...))
		return "added to pool"
	} else {
		//create an empty board for the game
		emptyBoard := GameBoard{'-', '-', '-', '-', '-', '-', '-', '-', '-'}
		//get playerO's address
		playerO := state.ReadBytesByKey("waitingPool")[:20]

		//update the waiting pool
		state.WriteBytesByKey("waitingPool", state.ReadBytesByKey("waitingPool")[20:])

		//create the new game
		newGame := Game{emptyBoard, callerAddress, playerO, true}
		//if there are free indexes, add a new game to one of them
		if len(state.ReadBytesByKey("freeIndexes")) != 0 {

			//get the free index
			ByteIndex := state.ReadBytesByKey("freeIndexes")[:8]
			index := bytesToUint64(ByteIndex)

			//remove the free index from the state
			state.WriteBytesByKey("freeIndexes", state.ReadBytesByKey("freeIndexes")[8:])

			games[index] = newGame

			//update the player's index:
			state.WriteUint64ByAddress(callerAddress, index)
			state.WriteUint64ByAddress(playerO, index)

		} else {
			//if there aren't any free indexes, append the new games to the games
			games = append(games, newGame)

			//update the player's index:
			state.WriteUint64ByAddress(callerAddress, uint64(len(games)-1))
			state.WriteUint64ByAddress(playerO, uint64(len(games)-1))
		}

	}
	//update the games to the state
	newGamesBytes, err := games.MarshalJSON()
	if err != nil {
		panic(err)
	}
	state.WriteBytesByKey("games", newGamesBytes)
	return "started new game"
}

func play(box uint32) (feedback string) {
	//get the games
	var games Games
	byteGames := state.ReadBytesByKey("games")
	err := games.UnmarshalJSON(byteGames)
	if err != nil {
		panic(err)
	}
	//call the method in order to play
	games.play(box)
	if checkIfWon() == 1 {
		return "congrats, you won"
	}
	return getBoard()

}
func (games *Games) play(box uint32) {

	signerAddress := address.GetSignerAddress()

	//get the relevant game
	game := (*games)[state.ReadUint64ByAddress(address.GetSignerAddress())]

	//if the player isn't registered for the game, don't let him play
	if !(bytes.Equal(signerAddress, game.PlayerX) || bytes.Equal(signerAddress, game.PlayerO)) {
		panic("you are not registered in any game game")
	}
	//make sure every player play's on his turn
	if game.PlayerXTurn {
		if bytes.Equal(signerAddress, game.PlayerO) {
			panic("it is player X's turn, wait for your turn")
		}
	} else {
		if bytes.Equal(signerAddress, game.PlayerX) {
			panic("it is player O's turn, wait for your turn")
		}
	}

	//make sure the index of the box is OK
	if 8 < box {
		panic("the index must be between 0 and 8")
	}
	//make sure the box isn't taken
	if game.Board[box] != '-' {
		panic("that index is taken")
	}

	//insert the rune according to the player
	if game.PlayerXTurn {
		game.Board[box] = 'X'
	} else {
		game.Board[box] = 'O'
	}
	//change the turn and update games
	game.PlayerXTurn = !game.PlayerXTurn
	(*games)[state.ReadUint64ByAddress(signerAddress)] = game

	//update state
	b, err := games.MarshalJSON()
	if err != nil {
		panic(err)
	}

	state.WriteBytesByKey("games", b)

}

func getBoard() (board string) {
	//get the games
	var games Games
	byteGames := state.ReadBytesByKey("games")
	err := games.UnmarshalJSON(byteGames)
	if err != nil {
		panic(err)
	}
	//get the relevant game
	signerAddress := address.GetSignerAddress()
	index := state.ReadUint64ByAddress(signerAddress)
	game := games[index]

	//if the player isn't registered for the game don't let him play
	if !(bytes.Equal(game.PlayerX, signerAddress) || bytes.Equal(game.PlayerO, signerAddress)) {
		panic("you are not registered in any game")
	}
	//return the board
	return string(game.Board)
}

func checkIfWon() (didWin uint32) {
	//get the games
	var games Games
	byteGames := state.ReadBytesByKey("games")
	err := games.UnmarshalJSON(byteGames)
	if err != nil {
		panic(err)
	}

	//call the helper methods and return the value accordingly
	if games.checkDiagonal() || games.checkColumn() || games.checkRow() {
		finishGame()
		return 1
	}
	return 0
}
func (games *Games) checkDiagonal() (didWin bool) {

	//get the game and the board
	signerAddress := address.GetSignerAddress()
	index := state.ReadUint64ByAddress(signerAddress)
	game := (*games)[index]
	board := game.Board

	//get the player sign
	var playerSign rune
	if bytes.Equal(signerAddress, game.PlayerX) {
		playerSign = 'X'
	} else if bytes.Equal(signerAddress, game.PlayerO) {
		playerSign = 'O'
	} else {
		panic("you are not registered for this game")
	}
	//check if any of the diagonal have 3 in a row
	if board[0] == playerSign && board[4] == playerSign && board[8] == playerSign {
		return true
	}
	if board[2] == playerSign && board[4] == playerSign && board[6] == playerSign {
		return true
	}
	return false

}
func (games *Games) checkColumn() (didWin bool) {
	//get the game and board
	signerAddress := address.GetSignerAddress()
	index := state.ReadUint64ByAddress(signerAddress)
	game := (*games)[index]
	board := game.Board
	//get the player's sign
	var playerSign rune
	if bytes.Equal(signerAddress, game.PlayerX) {
		playerSign = 'X'
	} else if bytes.Equal(signerAddress, game.PlayerO) {
		playerSign = 'O'
	} else {
		panic("you are not registered for this game")
	}

	//check if any of the columns have 3 in a row
	for i := 0; i < 3; i++ {
		for j := i; j <= i+6; j += 3 {
			if board[i+j] != playerSign {
				break
			}
			if j == i+6 {
				return true
			}
		}
	}
	return false
}
func (games *Games) checkRow() (didWin bool) {
	//get the game and board
	signerAddress := address.GetSignerAddress()
	index := state.ReadUint64ByAddress(signerAddress)
	game := (*games)[index]
	board := game.Board

	//get the player's sign
	var playerSign rune
	if bytes.Equal(signerAddress, game.PlayerX) {
		playerSign = 'X'
	} else if bytes.Equal(signerAddress, game.PlayerO) {
		playerSign = 'O'
	} else {
		panic("you are not registered for this game")
	}

	//check if any of the diagonal's have 3 in a row
	for i := 0; i <= 6; i += 3 {
		for j := i; j < i+3; j++ {
			if board[j] != playerSign {
				break
			}
			if j == i+2 {
				return true
			}
		}
	}
	return false
}

func finishGame() {
	//get the games
	signerAddress := address.GetSignerAddress()
	index := state.ReadUint64ByAddress(signerAddress)
	byteGames := state.ReadBytesByKey("games")
	var games Games
	err := games.UnmarshalJSON(byteGames)
	if err != nil {
		panic(err)
	}
	//get the relevant game
	game := games[index]
	//only a registered player can finish a game
	if !(bytes.Equal(signerAddress, game.PlayerX) || bytes.Equal(signerAddress, game.PlayerO)) {
		panic("you are not registered for this game")
	}
	//update the state
	game.PlayerO = []byte{}
	game.PlayerX = []byte{}
	games[index] = game
	b, err := games.MarshalJSON()
	if err != nil {
		panic(err)
	}
	state.WriteBytesByKey("games", b)
	state.WriteBytesByKey("freeIndexes", append(state.ReadBytesByKey("freeIndexes"), uint64ToBytes(index)...))
	state.ClearByAddress(game.PlayerO)
	state.ClearByAddress(game.PlayerX)
}

func quitPool() {
	waitingPool := state.ReadBytesByKey("waitingPool")
	inPool := false
	i := 0
	for ; i < len(waitingPool); i += 20 {
		if bytes.Equal(waitingPool[i:i+20], address.GetSignerAddress()) {
			inPool = true
			break
		}
	}
	if !inPool {
		panic("you are not in the pool")
	}

	waitingPool = append(waitingPool[:i], waitingPool[i+20:]...)
	state.WriteBytesByKey("waitingPool", waitingPool)
}

//custom types for the contract
type GameBoard []rune

type Game struct {
	Board       GameBoard `json:"board"`
	PlayerX     []byte    `json:"playerX"`
	PlayerO     []byte    `json:"playerO"`
	PlayerXTurn bool      `json:"playerXTurn"`
}
type Games []Game

//marshal and unmarshal methods for the types
func (board *GameBoard) MarshalJSON() (b []byte, err error) {
	if len(*board) != 9 {
		panic("invalid board size")
	}
	var s string
	s = string(*board)
	return json.Marshal(s)
}

func (board *GameBoard) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	for i := 0; i < 9; i++ {
		*board = append(*board, rune(s[i]))
	}

	return nil
}

func (game *Game) MarshalJSON() (b []byte, err error) {
	mappedGame := make(map[string][]byte)

	//write the board value to the map
	mappedGame["board"], err = game.Board.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}

	//write the player's to the map
	mappedGame["playerX"] = game.PlayerX
	mappedGame["playerO"] = game.PlayerO

	//write the player's turn to the map
	if game.PlayerXTurn {
		mappedGame["playerXTurn"] = []byte{1}
	} else {
		mappedGame["playerXTurn"] = []byte{0}
	}

	//return the marshaled map
	return json.Marshal(mappedGame)
}

func (game *Game) UnmarshalJSON(b []byte) (err error) {
	//unmarshal the data into a map
	mappedGame := make(map[string][]byte)
	err = json.Unmarshal(b, &mappedGame)
	if err != nil {
		panic(err)
	}

	boardBytes := mappedGame["board"]
	var board GameBoard
	err = board.UnmarshalJSON(boardBytes)
	if err != nil {
		return err
	}
	game.Board = board
	game.PlayerX = mappedGame["playerX"]
	game.PlayerO = mappedGame["playerO"]

	if mappedGame["playerXTurn"][0] == 1 {
		game.PlayerXTurn = true
	} else {
		game.PlayerXTurn = false
	}
	return nil
}

func (games *Games) MarshalJSON() (b []byte, err error) {
	gamesMap := make(map[int][]byte)
	//the key '0' is the length of the array
	gamesMap[0] = uint64ToBytes(uint64(len(*games)))
	for i := 1; i <= len(*games); i++ {
		gamesMap[i], err = (*games)[i-1].MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
	}
	return json.Marshal(gamesMap)

}

func (games *Games) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte{}) {
		return nil
	}
	var gamesMap map[int][]byte
	err = json.Unmarshal(b, &gamesMap)
	if err != nil {
		return err
	}
	var tempGame Game
	length := bytesToUint64(gamesMap[0])
	for i := 1; i <= int(length); i++ {
		err = tempGame.UnmarshalJSON(gamesMap[i])
		if err != nil {
			return err
		}
		*games = append(*games, tempGame)
	}

	return nil
}
