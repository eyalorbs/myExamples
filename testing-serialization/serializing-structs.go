package main

import (
	"encoding/json"
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
	game.PlayerX = mappedGame["playerX*"]
	game.PlayerO = mappedGame["playerO"]

	if mappedGame["playerXTurn"][0] == 1 {
		game.PlayerXTurn = true
	} else {
		game.PlayerXTurn = false
	}
	return nil
}

func main() {

	//emptyBoard := GameBoard{'-', '-', '-', '-', '-', '-', '-', '-', '-'}

	/*
		game1 := Game{emptyBoard, []byte{3,5,2,5}, []byte{4,2,6,240}, true}
		game2 := Game{emptyBoard, []byte{3,5,2,5,254}, []byte{4,2,6,240},  false}
	*/

}

type GameBoard []rune

type Game struct {
	Board       GameBoard `json:"board"`
	PlayerX     []byte    `json:"playerX"`
	PlayerO     []byte    `json:"playerO"`
	PlayerXTurn bool      `json:"playerXTurn"`
}
type Games []Game
