package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/state"
	"log"
)

func main() {
	player1 := []byte{42, 2, 5, 1, 5}
	player2 := []byte{26, 26, 2, 25}
	boardHashed := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}
	playerHits := uint8(0)
	defaultGuesses := guesses{}
	playerTurn := true
	defaultLastGuesses := coordinate{}
	boardApproved := uint8(3)

	thisGame := game{player1, player2, boardHashed, boardHashed, playerHits, playerHits, defaultGuesses, defaultGuesses, playerTurn, defaultLastGuesses, boardApproved, boardApproved}
	fmt.Println(thisGame)
	b, err := thisGame.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	var dec game
	err = dec.UnmarshalJSON(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dec)

	b, err = hex.DecodeString("7b2230223a22635262694275327544304a47652f647456756f46456b6c4d7942493d222c2231223a22683038636f3878644843465a3977596d776430643242694b2b35593d222c223130223a2241773d3d222c223131223a2241773d3d222c2232223a2247612f5a5a487265547151334757525849516275466a75686f785a6961546c76424a544a6a4a38483462633d222c2233223a2238344a6f374c2b326d347243665a5a374b5846627366626836424758676569426b6a64795178344b555a493d222c2234223a2245513d3d222c2235223a2242673d3d222c2236223a2265794977496a6f695a586c4a4e45394453545a4e553364705430527261553971526a6c4255543039496977694d534936496d56355354525051306b3254576c336155394561326c50616b59355156453950534973496a4577496a6f695a586c4a4e45394453545a4f513364705430527261553971546a6c4255543039496977694d5445694f694a6c65556b3054304e4a4e6b355464326c50524774705432704f4f554652505430694c4349784d694936496d56355354525051306b32546b4e336155394561326c50616c49355156453950534973496a457a496a6f695a586c4a4e45394453545a4f553364705430527261553971556a6c4255543039496977694d5451694f694a6c65556b3054304e4a4e6b357064326c5052477470543270534f554652505430694c4349784e534936496d56355354525051306b32546c4e336155394561326c50616c59355156453950534973496a4532496a6f695a586c4a4e45394453545a4f615864705430527261553971566a6c4255543039496977694d694936496d56355354525051306b3254586c336155394561326c50616b59355156453950534973496a4d694f694a6c65556b3054304e4a4e6b354464326c5052477470543270474f554652505430694c434930496a6f695a586c4a4e45394453545a4f553364705430527261553971526a6c4255543039496977694e534936496d56355354525051306b3254576c336155394561326c50616b6f355156453950534973496a59694f694a6c65556b3054304e4a4e6b313564326c50524774705432704b4f554652505430694c434933496a6f695a586c4a4e45394453545a4f513364705430527261553971536a6c4255543039496977694f434936496d56355354525051306b32546c4e336155394561326c50616b6f355156453950534973496a6b694f694a6c65556b3054304e4a4e6b313564326c50524774705432704f4f5546525054306966513d3d222c2237223a2265794977496a6f695a586c4a4e45394453545a4e553364705430527261553971526a6c4255543039496977694d534936496d56355354525051306b3254576c336155394561326c50616b6f355156453950534973496a4577496a6f695a586c4a4e45394453545a4e553364705430527261553971536a6c4251543039496977694d5445694f694a6c65556b3054304e4a4e6b317064326c50524774705432704f4f554642505430694c4349784d694936496d56355354525051306b3254586c336155394561326c50616c49355155453950534973496a457a496a6f695a586c4a4e45394453545a4f513364705430527261553971566a6c4251543039496977694d5451694f694a6c65556b3054304e4a4e6b395464326c50524774705432706f4f554642505430694c4349784e534936496d56355354525051306b3254586c336155394561326c50616b59355156453950534973496a4532496a6f695a586c4a4e45394453545a4f513364705430527261553971526a6b694c434979496a6f695a586c4a4e45394453545a4e655864705430527261553971546a6c4255543039496977694d794936496d56355354525051306b32546b4e336155394561326c50616c49355156453950534973496a51694f694a6c65556b3054304e4a4e6b355464326c5052477470543270574f554652505430694c434931496a6f695a586c4a4e45394453545a4f615864705430527261553971576a6c4251543039496977694e694936496d56355354525051306b32546e6c336155394561326c50616d51355155453950534973496a63694f694a6c65556b3054304e4a4e6b394464326c50524774705432706f4f554642505430694c434934496a6f695a586c4a4e45394453545a5055336470543052726155397162446c4251543039496977694f534936496d56355354525051306b325456524263306c715a7a564a616d39345455677751534a39222c2238223a2241413d3d222c2239223a22657949344f4349364e4377694f446b694f6a4639227d")
	if err != nil {
		log.Fatal(err)
	}
	err = dec.UnmarshalJSON(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dec)
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