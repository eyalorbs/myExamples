package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {

	//create and empty board

	var board Board
	for i := 0; i < 9; i++ {
		board[i] = '-'
	}
	for _, val := range board {
		fmt.Print(string(val))
	}
	fmt.Println("")

	b, err := json.Marshal(board)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b)

	var unmarshaledBoard Board

	err = json.Unmarshal(b, &unmarshaledBoard)
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range unmarshaledBoard {
		fmt.Print(string(val))
	}
	fmt.Println("")
	unmarshaledBoard.set(3)
	for _, val := range unmarshaledBoard {
		fmt.Print(string(val))
	}
}

type Board [9]rune

func (board *Board) set(index uint8) {
	board[index] = 'j'
}
