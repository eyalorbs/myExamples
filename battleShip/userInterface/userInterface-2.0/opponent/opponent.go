package main

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
	//set the reader
	reader := bufio.NewReader(os.Stdin)
	//deploy contract
	out := gammaCli.Run("deploy battleShip/backend/battleship.go")
	out = gammaCli.Run("deploy battleShip/backend/winnerContract/winnerContract.go")
	//get the ships
	boats := ships{}
	//user interface to get ships
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

	//auto generated ships
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

	//variables for the loop:

	//get the secret key
	fmt.Println("enter your super secret confidential secret key, don't tell anyone or you'll suffer from the devastating effects of losing to a cheater")
	secretKey, _ := reader.ReadString('\n')

	//get the hashed string
	b, marshaledShips, err := boats.sha256(secretKey)
	if err != nil {
		log.Fatal(err)
	}
	hashedShips := hex.EncodeToString(b)

	//get the points
	numOfHits := uint8(0)

	//get the user
	fmt.Println("\nwhat user are you?")
	var user int
	_, err = fmt.Scanf("%d", &user)
	if err != nil {
		log.Fatal("invalid input")

	}

	//get the player's ship coordinates
	coo := coordinates{}
	coo.getCoordinates(boats)

	//create the board
	board := Board{}
	board.new()

	//response
	var resp response

	for {
		//read the input from the user
		fmt.Println("\nenter your desired command")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		//handle input
		if input == "startGame" {
			//request to start game
			out = gammaCli.Run("send-tx battleShip/jsons/startGame.json -arg1 " + hashedShips + " -signer user" + strconv.Itoa(user))
			resp.getResponse(out)
			//print the event, go to the next iteration if found
			events, failed := resp.getEventsStartGame()
			fmt.Println(events)
			if failed {
				continue
			}

			//double break or continue if requested by inner loop

			printedMessage := false
			if events == "added to pool" {
				for { //print message only once
					if !printedMessage {
						fmt.Println("you will be asked to start to play once an opponent is found")
						printedMessage = true
					}
					//check if in game
					out = gammaCli.Run("send-tx battleShip/jsons/checkIfInGame.json -signer user" + strconv.Itoa(user))
					//get the response
					resp.getResponse(out)
					err, doContinue := resp.getEventCheckIfInGame()
					if err != nil {
						fmt.Println(err)
					}
					if doContinue {
						continue
					}
					sentMessage := false

					for {
						out = gammaCli.Run("run-query battleShip/jsons/getOpponentStatus.json -signer user" + strconv.Itoa(user))
						resp.getResponse(out)
						x, y, err, doBreak, doContinue := resp.getOpponentGuess()
						if err != nil {
							if err.Error() == "opponent didn't play yet" {
								if !sentMessage {
									fmt.Println(err)
									sentMessage = !sentMessage
								}
							} else {
								fmt.Println(err)
							}

						}
						if doContinue {
							continue
						}
						if doBreak {
							break
						}
						//check if the opponent's guess hit and update
						if coo.exists(x, y) {
							gammaCli.Run("send-tx battleShip/jsons/updateHit.json -arg1 1 -signer user" + strconv.Itoa(user))
							break
						} else {
							gammaCli.Run("send-tx battleShip/jsons/updateHit.json -arg1 0 -signer user" + strconv.Itoa(user))
							break
						}
					}
					break

				}
				//once done, continue to next itteration
				fmt.Println("started new game, it is your turn")
				continue
				//if a new game started continue
			} else if events == "started new game" {
				fmt.Println("started new game, it is your turn")
				continue
				//handle unexpected event
			} else {
				fmt.Println("unexpected event: " + events)
				continue
			}

		} else if input == "guess" {

			//print opponent's board
			fmt.Println("opponent's board:")
			board.printBoard()
			//get the guess:
			//get the coordinate
			x, y, err := getCoordinatesFromUser()
			if err != nil {
				fmt.Println("invalid input")
				continue
			}
			//run command and get response
			out = gammaCli.Run("send-tx battleShip/jsons/guess.json -arg1 " + strconv.Itoa(int(x)) + " -arg2 " + strconv.Itoa(int(y)) + " -signer user" + strconv.Itoa(user))
			resp.getResponse(out)
			event, panics, err := resp.getEventGuess()
			//if there is an error print it and continue to the next iteration
			if err != nil {
				fmt.Println(err)
				continue
			}
			doContinue := false
			if event == "approve your board" || panics == "both players need to approve their board" {
				out = gammaCli.Run("send-tx battleShip/jsons/approveBoard.json -arg1 " + secretKey + " -arg2 " + hex.EncodeToString(marshaledShips) + " -signer user" + strconv.Itoa(user))
				for {
					out = gammaCli.Run("send-tx battleship/jsons/finishGame.json -signer user" + strconv.Itoa(user))
					out = gammaCli.Run("send-tx battleShip/jsons/checkIfInGame.json -signer user" + strconv.Itoa(user))
					//get the response
					resp.getResponse(out)
					_, dobreak := resp.getEventCheckIfInGame()
					if dobreak {
						out = gammaCli.Run("send-tx battleShip/jsons/checkIfWon.json -signer user" + strconv.Itoa(user))
						resp.getResponse(out)
						if resp.OutputArguments[0].Value == "1" {
							log.Fatal("you lost... try not to lose next time... maybe that way you'll win")
						}
						if resp.OutputArguments[0].Value == "2" {
							log.Fatal("congratulations! you won")
						} else {
							fmt.Println("didn't think this would happen fkjdszk;rjdfwu54eijordsfndfjndzjncvx;jnxdl;kcmvdflz;kjvl;kxjvlk;fjv;lkxfjv;dfzjv;fjv;zfdjvoshrtoifjdvxnczjfxkfjdlkfjdlksfjldksfjlkdsjflksdjflkdsjflksjflkjdklfjldkjflk")
							continue
						}
					}
				}

			} else if event == "you already approved your board, we are in the endgame now" {
				fmt.Println(event)
				continue

			} else if event == "guess submitted" {
				fmt.Println("waiting for opponent to update hit...")
				for {
					out = gammaCli.Run("send-tx battleShip/jsons/didOpponentUpdateHit.json -signer user" + strconv.Itoa(user))
					resp.getResponse(out)
					didUpdate, err := resp.getEventDidOpponentUpdateHit()
					//if there is an error, exit this loop, and move to the next iteration on the next loop
					if err != nil {
						fmt.Println(err)
						doContinue = true
						break
					}
					//if opponent updated the hits, break
					if didUpdate {
						break
					}

				}
				//if told to continue, continue
				if doContinue {
					continue
				}
				//gets here after hits are updated
				out = gammaCli.Run("send-tx battleShip/jsons/getMyHits.json -signer user" + strconv.Itoa(user))
				resp.getResponse(out)
				newHits, err := resp.getMyHits()
				if err != nil {
					panic(err)

				}
				//if hit
				if numOfHits < newHits {
					//update the hits
					numOfHits = newHits
					//update the board
					board[x-1][y-1] = '*'
					//notify the player
					fmt.Println("you hit!")
					board.printBoard()
					//if has enough point to win, prove the board is ok, and check if won
					if numOfHits == 17 {
						out = gammaCli.Run("send-tx battleShip/jsons/approveBoard.json -arg1 " + secretKey + " -arg2 " + hex.EncodeToString(marshaledShips) + " -signer user" + strconv.Itoa(user))

						fmt.Println("waiting for opponent to approve board")
						for {
							out = gammaCli.Run("send-tx battleship/jsons/finishGame.json -signer user" + strconv.Itoa(user))
							out = gammaCli.Run("send-tx battleShip/jsons/checkIfInGame.json -signer user" + strconv.Itoa(user))
							//get the response
							resp.getResponse(out)
							_, dobreak := resp.getEventCheckIfInGame()
							if dobreak {
								out = gammaCli.Run("send-tx battleShip/jsons/checkIfWon.json -signer user" + strconv.Itoa(user))
								resp.getResponse(out)
								if resp.OutputArguments[0].Value == "1" {
									log.Fatal("you lost... try not to lose next time... maybe that way you'll win")
								}
								if resp.OutputArguments[0].Value == "2" {
									log.Fatal("congratulations! you won")
								} else {
									fmt.Println("didn't expect this asffdjdkjfekrjflkaerjfh;ahjferawhfoierahfdshfaehjfposehjfaehjfpaerhfowaefhpseafhasfhparsfhparghaprghpidsuhfgawrhzpofjfgjg;arjgpasej poershjfporjgpsej")
									continue
								}
							}

						}
					}
					//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
					///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

					//if missed
				} else {
					fmt.Println("you missed :-(")
					board[x-1][y-1] = 'O'
					board.printBoard()
				}

				fmt.Println("waiting for your opponent to guess...")
				sentMessage := false
				for {
					out = gammaCli.Run("run-query battleShip/jsons/getOpponentStatus.json -signer user" + strconv.Itoa(user))
					resp.getResponse(out)
					x, y, err, doBreak, doContinue := resp.getOpponentGuess()
					if err != nil {
						if err.Error() == "opponent didn't play yet" {
							if !sentMessage {
								fmt.Println(err)
								sentMessage = !sentMessage
							}
						} else {
							fmt.Println(err)
						}

					}
					if doContinue {
						continue
					}
					if doBreak {
						break
					}
					if coo.exists(x, y) {
						gammaCli.Run("send-tx battleShip/jsons/updateHit.json -arg1 1 -signer user" + strconv.Itoa(user))
						break
					} else {
						gammaCli.Run("send-tx battleShip/jsons/updateHit.json -arg1 0 -signer user" + strconv.Itoa(user))
						break
					}

				}

			} else {
				fmt.Println("unrecognized response")
				continue
			}

		} else if input == "quit" {
			out = gammaCli.Run("send-tx battleShip/jsons/quitGame.json -signer user" + strconv.Itoa(user))

		} else if input == "test" {
			out = gammaCli.Run("send-tx battleShip/jsons/getPlayersAddress.json -signer user" + strconv.Itoa(user))
			resp.getResponse(out)
			fmt.Println(resp.OutputArguments[0])

		} else if input == "didWinLastGame" {
			out = gammaCli.Run("run-query battleShip/jsons/checkIfWon.json -signer user" + strconv.Itoa(user))
			resp.getResponse(out)
			if resp.OutputArguments[0].Type == "string" {
				fmt.Println("you are playing a game right now")
			} else {
				if resp.OutputArguments[0].Value == "0" {
					fmt.Println("you do not have any previous games played")
				} else if resp.OutputArguments[0].Value == "1" {
					fmt.Println("you lost the last game")
				} else {
					fmt.Println("you won the last game")
				}
			}
		} else {
			fmt.Println("that command is not recognized")
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

func (coo *coordinates) exists(x, y uint8) (exists bool) {
	for _, val := range *coo {
		if val.X == x && val.Y == y {
			return true
		}
	}
	return false
}

func (resp *response) getResponse(out string) {
	err := json.Unmarshal([]byte(out), resp)
	if err != nil {
		panic("could not read from gamma")
	}
}

func (resp *response) getEventsStartGame() (event string, failed bool) {
	returns := resp.getReturns()
	if len(returns) != 0 {
		return returns[0].Value, true
	}
	returns = resp.getEvents()
	if len(returns) != 1 {
		return "unexpected number of events", true
	} else {
		return returns[0].Value, false
	}
}

func (resp *response) getEventCheckIfInGame() (err error, doContinue bool) {
	//get the returns
	returns := resp.getReturns()
	//if there is a panic message, return it
	if len(returns) != 0 {
		return errors.New(returns[0].Value), true
	}
	//get the events
	returns = resp.getEvents()
	//if there aren't 3 parameters for the events, return an error
	if len(returns) != 3 {
		return errors.New("unexpected number of events"), true

	} else {
		//if the player is still in pool, continue
		if returns[0].Value == "player still in pool" {
			return nil, true

			//if player is in game get the x and y and do not continue. only continue if the input is not valid
		} else if returns[0].Value == "player is in a game" {
			return nil, false
		} else {
			return errors.New("unrecognized event"), true
		}
	}

}

func (resp *response) getEventGuess() (event string, panic string, err error) {
	//get the returns
	returns := resp.getReturns()
	//if there is a return value, it is a panic. pass it as an error
	if len(returns) != 0 {
		if returns[0].Value == "both players need to approve their board" {
			return "", returns[0].Value, nil
		}
		return "", "", errors.New(returns[0].Value)
	}
	//get the events
	returns = resp.getEvents()
	if len(returns) == 1 {
		return returns[0].Value, "", nil
	} else if len(returns) == 2 {
		return returns[1].Value, "", nil
	} else {
		return "", "", errors.New("unexpected number of events")
	}

}

func (resp *response) getEventDidOpponentUpdateHit() (didUpdate bool, err error) {
	returns := resp.getReturns()
	if len(returns) != 0 {
		return false, errors.New(returns[0].Value)
	}
	returns = resp.getEvents()
	if len(returns) != 1 {
		return false, errors.New("unexpected number of events")
	}
	updated, err := strconv.Atoi(returns[0].Value)
	if err != nil {
		return false, err
	}
	if updated == 1 {
		return true, nil
	}
	return false, nil
}

func (resp *response) getMyHits() (hits uint8, err error) {
	val := resp.getReturns()
	if len(val) != 1 {
		return 101, errors.New("unexpected number of return values")
	}
	if val[0].Type == "string" {
		return 101, errors.New(val[0].Value)
	} else if val[0].Type == "uint32" {
		Hits, err := strconv.Atoi(val[0].Value)
		if err != nil {
			return 101, errors.New("failure converting string to int")
		}
		return uint8(Hits), nil

	} else {
		return 101, errors.New("unexpected type")
	}
}

func (resp *response) getCheckIfWon() (player1Won, player2Won bool, err error) {
	returns := resp.getReturns()
	if len(returns) == 1 {
		return false, false, errors.New(returns[0].Value)
	}
	if len(returns) == 2 {
		player1, err := strconv.Atoi(returns[0].Value)
		if err != nil {
			return false, false, err
		}
		player2, err := strconv.Atoi(returns[1].Value)
		if err != nil {
			return false, false, err
		}
		return player1 == 1, player2 == 1, err
	}
	return false, false, errors.New("unexpected length of return values")
}

func (resp *response) getOpponentGuess() (x, y uint8, err error, doBreak, doContinue bool) {
	returns := resp.getReturns()
	if len(returns) == 1 {
		if returns[0].Value == "it is your turn, you are not the one who is supposed to validate it" {
			return 10, 10, errors.New(returns[0].Value), true, false
		} else if returns[0].Value == "opponent didn't play yet" {
			return 10, 10, errors.New("opponent didn't play yet"), false, true
		} else {
			return 10, 10, errors.New(returns[0].Value), true, false

		}
	} else if len(returns) == 2 {
		X, err := strconv.Atoi(returns[0].Value)
		if err != nil {
			return 10, 10, errors.New("error converting to int"), true, false
		}
		x = uint8(X)
		Y, err := strconv.Atoi(returns[1].Value)
		if err != nil {
			return 10, 10, errors.New("error converting to int"), true, false
		}
		y = uint8(Y)
		return x, y, nil, false, false
	} else {
		return 10, 10, errors.New("unexpected number of return values"), true, false
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
