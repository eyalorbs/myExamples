package test

import (
	"github.com/orbs-network/orbs-contract-sdk/go/testing/gamma"
	"strings"
	"testing"
)

func TestStressTX(t *testing.T) {
	gammaCli := gamma.Cli().Start()
	defer gammaCli.Stop()

	for i := 0; i < 1500; i++ {
		message := "send-tx -env testnet42 -i ../jsons/add-10.json"
		//add to the count and check
		out := gammaCli.Run(message)
		if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
			t.Fatal(`"ExecutionResult": "SUCCESS"`)
		}
	}
	/*
		for i := 12; i < 13; i++ {

			argument := Argument{"Uint32", "10"}
			arguments := []Argument{argument}
			input := JSONinput{"MyCounter" + strconv.Itoa(i), "add", arguments}
			b, err := json.MarshalIndent(input, "", "")
			if err != nil {
				t.Fatal(err)
			}

			f, err := os.OpenFile("file.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
			if err != nil {
				t.Fatal(err)
			}
			// write to file, f.Write()
			_, err = f.Write(b)
			if err != nil {
				t.Fatal(err)
			}

			message := "send-tx -env testnet42 -i file.json"
			//add to the count and check
			out := gammaCli.Run(message)
			if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
				t.Fatal(`"ExecutionResult": "SUCCESS"`)
			}
			err = os.Remove("file.json")
			if err != nil{
				t.Fatal(err)
			}
		}
	*/
}

/*
type Argument struct{
	Type string
	Value string
}

type JSONinput struct {
	ContractName string
	MethodName   string
	Arguments    []Argument
}*/
