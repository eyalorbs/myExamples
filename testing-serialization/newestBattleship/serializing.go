package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	var coo coordinate
	coo.new(15, 3)
	Guesses := guesses{coo, coo, coo}
	guesses2 := append(Guesses, coordinate{4, 2})
	player1 := []byte{3, 26, 240, 1}
	player2 := []byte{84, 83, 84, 2}
	hashedBoard1 := []byte{63, 72, 72, 7, 2}
	hashedBoard2 := []byte{73, 72, 8, 2, 43}

	var match game
	match.new(player1, player2, hashedBoard1, hashedBoard2, 4, 0, Guesses, guesses2, false, coo)
	matches := make(games)

	fmt.Println(matches)

	b, err := matches.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	stuff := make(games)
	err = stuff.UnmarshalJSON(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stuff)

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
