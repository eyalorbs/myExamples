package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	coo1 := coordinate{3, 5}
	coo2 := coordinate{5, 2}
	box1 := box{coo1, true}
	box2 := box{coo2, false}
	box3 := box{coo1, false}
	box4 := box{coo2, true}
	boxes := secretBoard{box1, box2, box3, box4}
	fmt.Println("before being shuffled", boxes)
	boxes.Shuffle()
	fmt.Println("after being shuffled", boxes)
	b := boxes.sha256("password")
	fmt.Println(b)

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
func (boxes *secretBoard) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make(secretBoard, len(*boxes))
	perm := r.Perm(len(*boxes))
	for i, randIndex := range perm {
		ret[i] = (*boxes)[randIndex]
	}
	*boxes = ret
}
