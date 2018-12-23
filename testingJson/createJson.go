package main

import (
	"encoding/json"
	"log"
	"os"
)

type Argument struct {
	Type  string
	Value string
}

type JSONinput struct {
	ContractName string
	MethodName   string
	Arguments    []Argument
}

func main() {
	argument := Argument{"Uint32", "10"}
	arguments := []Argument{argument}
	input := JSONinput{"MyCounter", "add", arguments}

	b, err := json.MarshalIndent(input, "", "")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile("file.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal(err)
	}
	// write to file, f.Write()
	_, err = f.Write(b)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove("file.json")
	if err != nil {
		log.Fatal(err)
	}

	/*
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	*/
}
