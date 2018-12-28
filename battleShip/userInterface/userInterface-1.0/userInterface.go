package userInterface_1_0

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"math"
	"myExamples/silentGamma"
	"os"
	"strconv"
	"strings"
)

func main() {

	gammaCli := silentGamma.Cli().Start()
	defer gammaCli.Stop()
	reader := bufio.NewReader(os.Stdin)

	out := gammaCli.Run("deploy battleShip/backend/battleship.go")
	//get the ships
	boats := ships{}
	/*
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

	*/
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
	numOfHits := 0

	//saves the coordinates
	coo := coordinates{}
	coo.getCoordinates(boats)

	//get the secret key
	fmt.Println("enter your super secret confidential secret key, don't tell anyone or you'll suffer from the devastating effects of losing to a cheater")
	secretKey, _ := reader.ReadString('\n')

	//get the hashed string
	b, marshaledShips, err := boats.sha256(secretKey)
	if err != nil {
		log.Fatal(err)
	}
	hashedShips := hex.EncodeToString(b)

	fmt.Println("\nwhat user are you?")
	var user int
	_, err = fmt.Scanf("%d", &user)
	if err != nil {
		log.Fatal("invalid input")

	}
	fmt.Println("you are user", user)
	var resp response

	board := Board{}
	board.new()

	for {

		fmt.Println("\nenter your desired command")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		switch input {
		case "startGame":
			//start a game with the hashed ships
			out = gammaCli.Run("send-tx battleShip/jsons/startGame.json -arg1 " + hashedShips + " -signer user" + strconv.Itoa(user))
			err = json.Unmarshal([]byte(out), &resp)
			if err != nil {
				fmt.Println(err)
				continue
			}

			//if there are returns, it is a panic and therefor an error
			answer := resp.getReturns()
			if answer != nil {
				fmt.Println(answer[0].Value)
				continue
				//if it as event, print the event
			} else {
				answer = resp.getEvents()
				if answer == nil {
					fmt.Println("some error")
					continue
				}
				feedback := answer[0].Value
				if feedback == "added to pool" {
					fmt.Println("added to pool, you will be asked to enter a command when the game starts")
					for {
						out = gammaCli.Run("send-tx battleShip/jsons/checkIfInGame.json")
						err = json.Unmarshal([]byte(out), &resp)
						if err != nil {
							fmt.Println(err)
							continue
						}
						feedback = resp.getEvents()[0].Value
						if feedback == "player still in pool" {
							continue
						}
						if feedback == "player is in a game" {
							break
						}
						returns := resp.getReturns()
						if len(returns) == 0 {
							panic("unexpected non returns")
						}

						fmt.Println(returns[0].Value)
						break

					}

					for {
						out = gammaCli.Run("send-tx battleShip/jsons/getOpponentStatus.json -signer user" + strconv.Itoa(user))
						err = json.Unmarshal([]byte(out), &resp)
						if err != nil {
							panic(err)
						}
						x, y, err := checkIfOpponentPlayed(resp)
						if err != nil {
							continue
						}
					}
				} else if feedback == "started new game" {
					fmt.Println(feedback)
				}
			}

		case "guess":
			fmt.Println("this is your opponent's board: ")
			board.printBoard()
			//get the coordinate
			x, y, err := getCoordinatesFromUser()
			if err != nil {
				fmt.Println(err)
				continue
			}

			//run the command
			out = gammaCli.Run("send-tx battleShip/jsons/guess.json -arg1 " + strconv.Itoa(int(x)) + " -arg2 " + strconv.Itoa(int(y)) + " -signer user" + strconv.Itoa(user))
			err = json.Unmarshal([]byte(out), &resp)
			if err != nil {
				panic(err)
			}

			//if there are returns, it is a panic and therefor an error
			answer := resp.getReturns()
			if answer != nil {
				fmt.Println(answer[0].Value)
				continue

				//if it an event, print the event
			} else {
				answer = resp.getEvents()
				if answer == nil {
					fmt.Println(answer[0].Value)
					continue
				}
				event := answer[0].Value

				//if required, auto approve board
				if event == "you need to approve your board" {
					out = gammaCli.Run("send-tx battleShip/jsons/approveBoard.json -arg1 " + secretKey + " -arg2 " + hex.EncodeToString(marshaledShips) + " -signer user" + strconv.Itoa(user))
					fmt.Println("automatically approved your board")
					continue

					//if player approved board, wait for opponent to approve board
				} else if event == "you already approved your board, we are in the endgame now" {
					fmt.Println("wait for opponent to approve board")
					continue

					//if everything is ok
				} else if event == "guess submitted" {
					fmt.Println("guess submitted")
					//handle unknown event
				} else {
					panic("unexpected event:\n" + fmt.Sprint(event))
				}
				breaks := false
				count := 0
				for {

					out = gammaCli.Run("send-tx battleShip/jsons/didOpponentUpdateHit.json -signer user" + strconv.Itoa(user))
					err = json.Unmarshal([]byte(out), &resp)
					if err != nil {
						panic("the error is here: " + fmt.Sprint(err))
					}
					feedback := resp.getReturns()
					if len(feedback) != 0 {
						fmt.Println(feedback[0].Value)
						breaks = true
						break
					}
					feedback = resp.getEvents()
					if len(feedback) != 1 {
						panic("unexpected number of events")
					}
					ready, err := strconv.Atoi(feedback[0].Value)
					if err != nil {
						panic("unexpected error")
					}
					if ready == 1 {
						break
					}
					if count == 0 {
						fmt.Println("waiting for your opponent to respond to your guess...")
						count++
					}
				}
				if breaks {
					break
				}
				//get the number of hits
				out = gammaCli.Run("send-tx battleShip/jsons/getMyHits.json")
				err = json.Unmarshal([]byte(out), &resp)
				if err != nil {
					panic(err)
				}
				//if there are unexpected events, panic
				events := resp.getEvents()
				if len(events) != 0 {
					panic(events[0].Value)
				}
				//if there is more than one return value, panic
				returns := resp.getReturns()
				if len(returns) != 1 {
					panic("unexpected length of return values:\n" + fmt.Sprint(returns))
				}
				//convert the number form string to int
				hits, err := strconv.Atoi(returns[0].Value)
				if err != nil {
					panic(err)
				}
				//notify the player if he hit or missed
				if numOfHits < hits {
					fmt.Println("you hit!")
					board.addHit(x, y)
					numOfHits = hits
				} else {
					fmt.Println("you missed ):")
					board.addMiss(x, y)
				}

				//check if player won, if so, approve his board and check if he won
				if numOfHits == 17 {
					out = gammaCli.Run("send-tx battleShip/jsons/approveBoard.json -arg1 " + secretKey + " -arg2 " + hex.EncodeToString(marshaledShips) + " -signer user" + strconv.Itoa(user))
					err = json.Unmarshal([]byte(out), &resp)
					if err != nil {
						panic(err)
					}
					returns = resp.getReturns()
					if len(returns) != 0 {
						panic(returns)
					}

					out = gammaCli.Run("send-tx battleShip/jsons/checkIfWon.json")
					err = json.Unmarshal([]byte(out), &resp)
					if err != nil {
						panic(err)
					}
					returns = resp.getReturns()
					if len(returns) != 0 {
						panic(returns)
					}

					fmt.Println("you have sunk all the ships and your board has been approved")
				}
				//wait for opponent to guess
				for {
					out = gammaCli.Run("send-tx battleShip/jsons/getOpponentStatus.json -signer user" + strconv.Itoa(user))
					guess := resp.getReturns()
					if len(guess) == 1 {
						if guess[0].Value == "both players need to approve their board" {
							break
						}
					}
					if len(guess) != 2 {
						continue
					}
					for _, val := range coo {
						x, err := strconv.Atoi(guess[0].Value)
						if err != nil {
							panic(err)
						}
						y, err := strconv.Atoi(guess[1].Value)
						if err != nil {
							panic(err)
						}
						if uint8(x) == val.X && uint8(y) == val.Y {
							out = gammaCli.Run("send-tx battleShip/jsons/updateHit.json -arg1 1")
							err = json.Unmarshal([]byte(out), &resp)

						} else {
							out = gammaCli.Run("send-tx battleShip/jsons/updateHit.json -arg1 0")
							err = json.Unmarshal([]byte(out), &resp)
						}
						if err != nil {
							panic(err)
						}
						returns := resp.getReturns()
						if len(returns) != 0 {
							panic(returns[0])
						}

					}
					break
				}
			}

		case "quitGame":
			out = gammaCli.Run("send-tx battleShip/jsons/quitGame.json")
			err = json.Unmarshal([]byte(out), &resp)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("successfully quit")
			break

		default:
			fmt.Println("command  not found")
		}

	}

}

func getCoordinatesFromUser() (x, y uint8, err error) {
	fmt.Println("enter the coordinates")
	_, err = fmt.Scanf("%d", &x)
	if err != nil {
		return 10, 10, err
	}
	_, err = fmt.Scanf("%d", &y)
	if err != nil {
		return 10, 10, err
	}
	return x, y, nil
}
func checkIfOpponentPlayed(resp response) (x, y uint8, err error) {
	returnVal := resp.getReturns()
	if len(returnVal) == 1 {
		return 10, 10, errors.New(returnVal[0].Value)
	} else if len(returnVal) == 2 {
		x, err := strconv.Atoi(returnVal[0].Value)
		if err != nil {
			return 10, 10, err
		}
		y, err := strconv.Atoi(returnVal[1].Value)
		if err != nil {
			return 10, 10, err
		}
		return uint8(x), uint8(y), nil
	} else {
		panic("unexpected number of returns")
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
	ContractName string
	EventName    string
	Arguments    []argument
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
		for j := 0; j < len(resp.OutputEvents[i].Arguments); j++ {
			values = append(values, resp.OutputEvents[i].Arguments[j])
		}

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
func (boats *ships) sha256(sk string) (sha []byte, marshaledShips []byte, err error) {
	h := hmac.New(sha256.New, []byte(sk))
	b, err := boats.MarshalJSON()
	if err != nil {
		return []byte{}, []byte{}, err
	}
	h.Write(b)
	return h.Sum(nil), b, nil

}

type Board [10][10]rune

func (board *Board) new() {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			board[i][j] = '-'
		}
	}
}
func (board *Board) addHit(x, y uint8) {
	board[x][y] = '*'
}
func (board *Board) addMiss(x, y uint8) {
	board[x][y] = 'O'
}
func (board *Board) printBoard() {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Print(string(board[i][j]))
		}
		fmt.Println("")
	}
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
