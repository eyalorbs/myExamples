package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	var game Game
	game.new([]byte{42, 24, 243, 5}, []byte{5, 26, 2, 6}, []byte{6, 26, 6, 62, 66, 45, 43}, []byte{63, 63, 62, 2})
	var games = Games{}

	fmt.Println(games)

	b, err := games.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	dec := make(Games)
	err = dec.UnmarshalJSON(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dec)
	dec["he"] = game
	fmt.Println(dec)
}

type Game struct {
	Player1     []byte
	Player2     []byte
	Player1Root []byte
	Player2Root []byte
}

func (game *Game) new(player1, player2, player1root, player2root []byte) {
	game.Player1 = player1
	game.Player2 = player2
	game.Player1Root = player1root
	game.Player2Root = player2root
}
func (game *Game) MarshalJSON() (b []byte, err error) {
	gameMap := make(map[byte][]byte)
	gameMap[0] = game.Player1
	gameMap[1] = game.Player2
	gameMap[2] = game.Player1Root
	gameMap[3] = game.Player2Root
	return json.Marshal(gameMap)
}

func (game *Game) UnmarshalJSON(b []byte) (err error) {
	gameMap := make(map[byte][]byte)
	err = json.Unmarshal(b, &gameMap)
	game.Player1 = gameMap[0]
	game.Player2 = gameMap[1]
	game.Player1Root = gameMap[2]
	game.Player2Root = gameMap[3]
	return nil
}

type Games map[string]Game

func (games *Games) MarshalJSON() (b []byte, err error) {
	gamesMap := make(map[string][]byte)
	for key, value := range *games {
		gamesMap[key], err = value.MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
	}
	return json.Marshal(gamesMap)
}

func (games *Games) UnmarshalJSON(b []byte) (err error) {
	var gamesMap = make(map[string][]byte)
	err = json.Unmarshal(b, &gamesMap)
	if err != nil {
		return err
	}
	var tempGame Game
	for key, value := range gamesMap {
		err = json.Unmarshal(value, &tempGame)
		if err != nil {
			return err
		}
		(*games)[key] = tempGame
	}
	return nil
}
