package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"myExamples/silentGamma"
	"os"
	"strings"
)

func main() {

	gammaCli := silentGamma.Cli().Start()
	defer gammaCli.Stop()
	reader := bufio.NewReader(os.Stdin)
	out := gammaCli.Run("deploy -name tictactoe -code tictactoe/tictactoe.go")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		log.Fatal("contract cannot be deployed")
	}
	fmt.Println("contract has been deployed")
	for {
		input, _ := reader.ReadString('\n')

		switch input {

		case "startGame\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSuffix(text, "\n")
			out = gammaCli.Run("send-tx -i tictactoe/jsons/startGame.json -signer user" + text)
			fmt.Println(getValue(out))

		case "play-0\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/play-0.json -signer user" + text)
			fmt.Println(getValue(out))

		case "play-1\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/play-1.json -signer user" + text)
			fmt.Println(getValue(out))

		case "play-2\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/play-2.json -signer user" + text)
			fmt.Println(getValue(out))

		case "play-3\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/play-3.json -signer user" + text)
			fmt.Println(getValue(out))

		case "play-4\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/play-4.json -signer user" + text)
			fmt.Println(getValue(out))

		case "play-5\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/play-5.json -signer user" + text)
			fmt.Println(getValue(out))

		case "play-6\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/play-6.json -signer user" + text)
			fmt.Println(getValue(out))

		case "play-7\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/play-7.json -signer user" + text)
			fmt.Println(getValue(out))

		case "play-8\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/play-8.json -signer user" + text)
			fmt.Println(getValue(out))

		case "quitPool\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/quitPool.json -signer user" + text)
			fmt.Println(getValue(out))

		case "finishGame\n":
			fmt.Println("which user are you?")
			text, _ := reader.ReadString('\n')
			out = gammaCli.Run("send-tx -i tictactoe/jsons/finishGame.json -signer user" + text)
			fmt.Println(getValue(out))
		default:
			fmt.Println("invalid input")
		}
	}
}

type outputArgument struct {
	Type  string
	Value string
}

type response struct {
	RequestStatus     string
	TxId              string
	ExecutionResult   string
	OutputArguments   []outputArgument
	TransactionStatus string
	BlockHeight       string
	BlockTimestamp    string
}

func getValue(out string) (value string) {
	var resp response
	b := []byte(out)
	err := json.Unmarshal([]byte(b), &resp)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.OutputArguments) == 0 {
		return ""
	}
	return resp.OutputArguments[0].Value
}
