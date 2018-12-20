package test

import (
	"github.com/orbs-network/orbs-contract-sdk/go/testing/gamma"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	gammaCli := gamma.Cli().Start()
	defer gammaCli.Stop()

	//check if the contract was successfully deployed
	out := gammaCli.Run("deploy -name tictactoe -code ../tictactoe.go")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("deploy failed")
	}

	//add user1 to the pool
	out = gammaCli.Run("send-tx -i ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}
	//check if user1 successfully quit pool
	out = gammaCli.Run("send-tx -i ../jsons/quitPool.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 could not quit pool")
	}

	//add user1 to the pool
	out = gammaCli.Run("send-tx -i ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}

	//add user2 to the game
	out = gammaCli.Run("send-tx -i ../jsons/startGame.json -signer user2")
	if !strings.Contains(out, `"Value": "started new game"`) {
		t.Fatal("adding user2 failed")
	}

	//check if the initial boar is successful
	out = gammaCli.Run("read -i ../jsons/getBoard.json -signer user1")
	if !strings.Contains(out, `"Value": "---------"`) {
		t.Fatal("initial board failed")
	}

	//check if user2 can successfully end the game
	out = gammaCli.Run("send-tx -i ../jsons/finishGame.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 could not end the game")
	}

	//add user1 to the pool
	out = gammaCli.Run("send-tx -i ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}

	//add user2 to the game
	out = gammaCli.Run("send-tx -i ../jsons/startGame.json -signer user2")
	if !strings.Contains(out, `"Value": "started new game"`) {
		t.Fatal("adding user2 failed")
	}

	//user1 should be stopped from playing first move
	out = gammaCli.Run("send-tx -i ../jsons/play-0.json -signer user1")
	if !strings.Contains(out, `"Value": "it is player X's turn, wait for your turn"`) {
		t.Fatal("user1 was not blocked")
	}
	//user2 should play
	out = gammaCli.Run("send-tx -i ../jsons/play-0.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user2 failed to play")
	}

	//check if the board was successfully updated
	out = gammaCli.Run("read -i ../jsons/getBoard.json")
	if !strings.Contains(out, `"Value": "X--------"`) {
		t.Fatal("updating board failed")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx -i ../jsons/play-0.json -signer user1")
	if !strings.Contains(out, `"Value": "that index is taken"`) {
		t.Fatal("user1 was not blocked")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx -i ../jsons/play-1.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 failed to play")
	}

	//check if the board was successfully updated
	out = gammaCli.Run("read -i ../jsons/getBoard.json")
	if !strings.Contains(out, `"Value": "XO-------"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx -i ../jsons/play-3.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user2 failed to play")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx -i ../jsons/play-4.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 failed to play")
	}

	//check if the board was successfully updated
	out = gammaCli.Run("read -i ../jsons/getBoard.json")
	if !strings.Contains(out, `"Value": "XO-XO----"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx -i ../jsons/play-6.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user2 failed to play")
	}
	//check if the board was successfully updated
	out = gammaCli.Run("read -i ../jsons/getBoard.json")
	if !strings.Contains(out, `"Value": "XO-XO-X--"`) {
		t.Fatal("updating board failed")
	}

	//check if won
	out = gammaCli.Run("send-tx -i ../jsons/checkIfWon.json -signer user2")
	if !strings.Contains(out, `"Value": "1"`) {
		t.Fatal("check If won failed")
	}

	//the game should have ended, lets start a new one
	//add user1 to the pool
	out = gammaCli.Run("send-tx -i ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}

	//add user2 to the game
	out = gammaCli.Run("send-tx -i ../jsons/startGame.json -signer user2")
	if !strings.Contains(out, `"Value": "started new game"`) {
		t.Fatal("adding user2 failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx -i ../jsons/play-0.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user2 failed to play")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx -i ../jsons/play-1.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 failed to play")
	}

	//check if the board was successfully updated
	out = gammaCli.Run("read -i ../jsons/getBoard.json")
	if !strings.Contains(out, `"Value": "XO-------"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx -i ../jsons/play-4.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user2 failed to play")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx -i ../jsons/play-7.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 failed to play")
	}

	//check if the board was successfully updated
	out = gammaCli.Run("read -i ../jsons/getBoard.json")
	if !strings.Contains(out, `"Value": "XO--X--O-"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx -i ../jsons/play-8.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user2 failed to play")
	}

	//check if game won
	out = gammaCli.Run("send-tx -i ../jsons/checkIfWon.json -signer user2")
	if !strings.Contains(out, `"Value": "1"`) {
		t.Fatal("check If won failed")
	}

	//the game should have ended, lets start a new one
	//add user1 to the pool
	out = gammaCli.Run("send-tx -i ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}

	//add user2 to the game
	out = gammaCli.Run("send-tx -i ../jsons/startGame.json -signer user2")
	if !strings.Contains(out, `"Value": "started new game"`) {
		t.Fatal("adding user2 failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx -i ../jsons/play-0.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user2 failed to play")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx -i ../jsons/play-5.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 failed to play")
	}

	//check if the board was successfully updated
	out = gammaCli.Run("read -i ../jsons/getBoard.json")
	if !strings.Contains(out, `"Value": "X----O---"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx -i ../jsons/play-1.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user2 failed to play")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx -i ../jsons/play-7.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 failed to play")
	}

	//check if the board was successfully updated
	out = gammaCli.Run("read -i ../jsons/getBoard.json")
	if !strings.Contains(out, `"Value": "XX---O-O-"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx -i ../jsons/play-2.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user2 failed to play")
	}

	//check if game won
	out = gammaCli.Run("send-tx -i ../jsons/checkIfWon.json -signer user2")
	if !strings.Contains(out, `"Value": "1"`) {
		t.Fatal("check If won failed")
	}

}
