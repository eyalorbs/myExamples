package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
	"math"
)

func shipsOk(boats ships) (ok bool, shipCoordinates []coordinate) {
	//if there aren't 5 boars return false
	if len(boats) != 5 {
		return false, nil
	}
	//go over each boat
	for _, val := range boats {
		//check if boat is diagonal
		fmt.Println(val.headCoordinates.Y, val.tailCoordinates.Y)
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
		if uint8(math.Abs(float64(val.headCoordinates.X)-float64(val.tailCoordinates.X)))+1 != length && uint8(math.Abs(float64(val.headCoordinates.Y)-float64(val.tailCoordinates.Y)))+1 != length {
			fmt.Println(val.tailCoordinates.X - val.headCoordinates.X)
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

func main() {
	boats := ships{}
	var boat1 ship
	var boat2 ship
	var boat3 ship
	var boat4 ship
	var boat5 ship
	boat1.new("Carrier", 1, 1, 5, 1)
	boat2.new("Battleship", 2, 2, 5, 2)
	boat3.new("Cruiser", 3, 3, 5, 3)
	boat4.new("Submarine", 4, 4, 6, 4)
	boat5.new("Destroyer", 5, 5, 6, 5)
	boats = ships{boat1, boat2, boat3, boat4, boat5}

	ok, _ := shipsOk(boats)
	fmt.Println(ok)
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
