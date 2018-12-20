package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

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

func uint64ToBytes(num uint64) (b []byte) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	return b
}

func bytesToUint64(b []byte) (num uint64) {

	return binary.LittleEndian.Uint64(b)

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

func main() {

	emptyBoard1 := GameBoard{'-', '-', '-', '-', '-', '-', '-', '-', '-'}
	emptyBoard2 := GameBoard{'-', '-', '-', '-', '-', '-', '-', '-', '-'}

	game1 := Game{emptyBoard1, []byte{3, 5, 2, 5}, []byte{4, 2, 6, 240}, true}
	game2 := Game{emptyBoard2, []byte{3, 5, 2, 5, 254}, []byte{4, 2, 6, 240}, false}

	games := Games{game1, game2}
	fmt.Println(games)

	games.play()
	fmt.Println(games[0].Board)
	fmt.Println(games[1].Board)

}

func (games *Games) play() {
	game := (*games)[0]
	game.Board[4] = 'X'
}

type GameBoard []rune

type Game struct {
	Board       GameBoard `json:"board"`
	PlayerX     []byte    `json:"playerX"`
	PlayerO     []byte    `json:"playerO"`
	PlayerXTurn bool      `json:"playerXTurn"`
}
type Games []Game
