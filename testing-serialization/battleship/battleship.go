package main

import (
	"encoding/json"
	"fmt"
	"github.com/orbs-network/orbs-contract-sdk/go/testing/gamma"
	"log"
	"strings"
)

func main() {
	gammaCli := gamma.Cli().Start()
	defer gammaCli.Stop()

	out := gammaCli.Run("deploy testing-serialization/battleship/simpleContract/simpleContract.go")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		log.Fatal("contract cannot be deployed")
	}

	out = gammaCli.Run("run-query testing-serialization/battleship/simpleContract/returnAndEvent.json")
	var resp response
	err := json.Unmarshal([]byte(out), &resp)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(resp.OutputArguments); i++ {
		fmt.Println("the return value is: ", resp.OutputArguments[i].Value)
	}
	for i := 0; i < len(resp.OutputEvents); i++ {
		fmt.Println("the contract is: ", resp.OutputEvents[i].ContractName, " the name of the event is: ", resp.OutputEvents[i].EventName, " and the argument is: ", resp.OutputEvents[i].eventArguments[i].Value)
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
