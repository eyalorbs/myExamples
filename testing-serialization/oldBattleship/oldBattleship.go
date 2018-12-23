package oldBattleship

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"math"
)

func main() {
	var head coordinate
	head.new(9, 5)
	var tail coordinate
	tail.new(9, 1)

	var boat ship
	err := boat.new(head, tail, "Carrier", 5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(boat)
	b, err := boat.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	var dec ship

	err = dec.UnmarshalJSON(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dec)
}

type coordinate struct {
	X uint32
	Y uint32
}

func (coo *coordinate) new(x, y uint32) {
	*coo = coordinate{x, y}
}

func (coo *coordinate) MarshalJSON() (b []byte, err error) {
	coordinateMap := make(map[rune][]byte)
	coordinateMap['X'] = uint32ToBytes(coo.X)
	coordinateMap['Y'] = uint32ToBytes(coo.Y)
	return json.Marshal(coordinateMap)
}

func (coo *coordinate) UnmarshalJSON(b []byte) (err error) {
	var coordinateMap map[rune][]byte
	err = json.Unmarshal(b, &coordinateMap)
	if err != nil {
		log.Fatal(err)
	}
	*coo = coordinate{bytesToUint32(coordinateMap['X']), bytesToUint32(coordinateMap['Y'])}
	return nil
}

type Coordinates []coordinate

func (coo *Coordinates) new(head, tail coordinate, shipSize int) (err error) {
	//check if in bounds
	if 9 < head.X || 9 < head.Y || 9 < tail.X || 9 < tail.Y {
		return errors.New("coordinates not in bounds")
	}
	//chick if diagonal
	if head.X != tail.X && head.Y != tail.Y {
		return errors.New("ship cannot be diagonal")
	}

	//create the slice
	for i := uint32(math.Min(float64(head.X), float64(tail.X))); i <= uint32(math.Max(float64(head.X), float64(tail.X))); i++ {
		for j := uint32(math.Min(float64(head.Y), float64(tail.Y))); j <= uint32(math.Max(float64(head.Y), float64(tail.Y))); j++ {
			*coo = append(*coo, coordinate{i, j})
		}
	}
	//check if length is valid
	if len(*coo) != shipSize {
		return errors.New("shipSize does not match desired ship size, the ship's size is: " + fmt.Sprint(len(*coo)))
	}

	return nil
}

func (coo *Coordinates) MarshalJSON() (b []byte, err error) {
	coordinatesMap := make(map[int][]byte)
	coordinatesMap[0] = uint32ToBytes(uint32(len(*coo)))

	for i := 1; i <= len(*coo); i++ {
		coordinatesMap[i], err = (*coo)[i-1].MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
	}
	return json.Marshal(coordinatesMap)

}

func (coo *Coordinates) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte{}) {
		return nil
	}
	var coordinatesMap map[int][]byte
	err = json.Unmarshal(b, &coordinatesMap)
	if err != nil {
		return err
	}
	var tempCoordinate coordinate
	length := bytesToUint64(coordinatesMap[0])
	for i := 1; i <= int(length); i++ {
		err = tempCoordinate.UnmarshalJSON(coordinatesMap[i])
		if err != nil {
			return err
		}
		*coo = append(*coo, tempCoordinate)
	}

	return nil
}

type ship struct {
	Coordinates Coordinates
	ShipName    string
	NumOfHits   uint32
}

func (boat *ship) new(headCoordinates, tailCoordinates coordinate, shipName string, shipSize int) (err error) {
	err = boat.Coordinates.new(headCoordinates, tailCoordinates, shipSize)
	if err != nil {
		return err
	}
	boat.ShipName = shipName
	boat.NumOfHits = 0

	return nil

}
func (boat *ship) MarshalJSON() (b []byte, err error) {
	shipMap := make(map[string][]byte)
	shipMap["Coordinates"], err = boat.Coordinates.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}
	shipMap["shipName"] = []byte(boat.ShipName)
	shipMap["numOfHits"] = uint32ToBytes(boat.NumOfHits)

	return json.Marshal(shipMap)
}

func (boat *ship) UnmarshalJSON(b []byte) (err error) {
	shipMap := make(map[string][]byte)
	err = json.Unmarshal(b, &shipMap)
	if err != nil {
		return err
	}

	var slice Coordinates
	err = slice.UnmarshalJSON(shipMap["Coordinates"])
	if err != nil {
		return err
	}

	boat.Coordinates = slice
	boat.ShipName = string(shipMap["shipName"])
	boat.NumOfHits = bytesToUint32(shipMap["numOfHits"])
	return nil

}

type playerBoard struct {
	Board      [10][10]rune
	Carrier    ship
	Battleship ship
	Cruiser    ship
	submarine  ship
	destroyer  ship
}

type game struct {
	player1      []byte
	player2      []byte
	player1Board playerBoard
	player2Board playerBoard
}

type games []game

func uint64ToBytes(num uint64) (b []byte) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	return b
}

func bytesToUint64(b []byte) (num uint64) {

	return binary.LittleEndian.Uint64(b)

}

func uint32ToBytes(num uint32) (b []byte) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint32(b, num)
	return b
}

func bytesToUint32(b []byte) (num uint32) {

	return binary.LittleEndian.Uint32(b)

}

//don't copy anything lower than here

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func New(text string) error {
	return &errorString{text}
}
