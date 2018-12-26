package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/orbs-network/orbs-contract-sdk/go/testing/gamma"
	"log"
	"math"
	"os"
)

func main() {

	gammaCli := gamma.Cli().Start()
	defer gammaCli.Stop()
	reader := bufio.NewReader(os.Stdin)
	//get the ships
	boats := ships{}
	for i := 0; i < 5; i++ {
		var boat ship
		fmt.Println("enter the name of the ship")
		name, _ := reader.ReadString('\n')
		fmt.Println("enter headX, headY, tailX, tailY")
		var headX uint8
		var headY uint8
		var tailX uint8
		var tailY uint8
		_, err := fmt.Scanf("%d", &headX)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Scanf("%d", &headY)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Scanf("%d", &tailX)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Scanf("%d", &tailY)
		if err != nil {
			log.Fatal(err)
		}

		boat.new(name, headX, headY, tailX, tailY)
		boats = append(boats, boat)
	}

	//saves the coordinates
	coo := coordinates{}
	coo.getCoordinates(boats)

	//get the secret key
	fmt.Println("enter your super secret confidential secret key, don't tell anyone or you'll suffer from the devastating effects of losing to a cheater")
	secretKey, _ := reader.ReadString('\n')

	//get the hashed string
	b, err := boats.sha256(secretKey)
	if err != nil {
		log.Fatal(err)
	}
	hashedShips := hex.EncodeToString(b)

	//print the str just so that the editor will shut up
	fmt.Println(hashedShips)

	var resp response

	for {
		fmt.Println("enter your desired command")
		input, _ := reader.ReadString('\n')
		switch input {
		case "startGame":
			out := gammaCli.Run("deploy battleShip/jsons/startGame.json -arg1 " + hashedShips)
			err = json.Unmarshal([]byte(out), resp)
			if err != nil {
				panic(err)
			}
			answer := resp.getReturns()
			if answer != nil {
				panic(answer[0].Value)
			} else {
				answer = resp.getEvents()
			}

		case "getOpponentStatus":

		case "guess":

		case "updateHit":

		case "quitGame":

		case "getMyHits":

		}
	}
}
func (coo *coordinates) getCoordinates(boats ships) {
	for _, val := range boats {
		for i := uint8(math.Min(float64(val.headCoordinates.X), float64(val.tailCoordinates.X))); i <= uint8(math.Max(float64(val.headCoordinates.X), float64(val.tailCoordinates.X))); i++ {
			for j := uint8(math.Min(float64(val.headCoordinates.Y), float64(val.tailCoordinates.Y))); j <= uint8(math.Max(float64(val.headCoordinates.Y), float64(val.tailCoordinates.Y))); j++ {
				tempCoo := coordinate{i, j}
				*coo = append(*coo, tempCoo)
			}
		}

	}
}

type argument struct {
	Type  string
	Value string
}
type outputEvent struct {
	ContractName   string
	EventName      string
	eventArguments []argument
}

type response struct {
	RequestStatus     string
	TxId              string
	ExecutionResult   string
	OutputArguments   []argument
	OutputEvents      []outputEvent
	TransactionStatus string
	BlockHeight       string
	BlockTimestamp    string
}

func (resp *response) getReturns() (values []argument) {

	for i := 0; i < len(resp.OutputArguments); i++ {
		values = append(values, resp.OutputArguments[i])
	}
	return values
}
func (resp *response) getEvents() (values []argument) {

	for i := 0; i < len(resp.OutputEvents); i++ {
		values = append(values, resp.OutputArguments[i])
	}
	return values
}

type coordinate struct {
	X uint8
	Y uint8
}
type coordinates []coordinate

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
