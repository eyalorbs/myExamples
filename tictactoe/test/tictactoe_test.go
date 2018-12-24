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
	out := gammaCli.Run("deploy ../tictactoe.go")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("deploy failed")
	}

	//add user1 to the pool
	out = gammaCli.Run("send-tx ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}
	//check if user1 successfully quit pool
	out = gammaCli.Run("send-tx ../jsons/quitPool.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 could not quit pool")
	}

	//add user1 to the pool
	out = gammaCli.Run("send-tx ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}

	//add user2 to the game
	out = gammaCli.Run("send-tx ../jsons/startGame.json -signer user2")
	if !strings.Contains(out, `"Value": "started new game"`) {
		t.Fatal("adding user2 failed")
	}

	//check if the initial boar is successful
	out = gammaCli.Run("run-query ../jsons/getBoard.json -signer user1")
	if !strings.Contains(out, `"Value": "---------"`) {
		t.Fatal("initial board failed")
	}

	//check if user2 can successfully end the game
	out = gammaCli.Run("send-tx ../jsons/finishGame.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("user1 could not end the game")
	}

	//add user1 to the pool
	out = gammaCli.Run("send-tx ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}

	//add user2 to the game
	out = gammaCli.Run("send-tx ../jsons/startGame.json -signer user2")
	if !strings.Contains(out, `"Value": "started new game"`) {
		t.Fatal("adding user2 failed")
	}

	//user1 should be stopped from playing first move
	out = gammaCli.Run("send-tx ../jsons/play-0.json -signer user1")
	if !strings.Contains(out, `"Value": "it is player X's turn, wait for your turn"`) {
		t.Fatal("user1 was not blocked")
	}
	//user2 should play
	out = gammaCli.Run("send-tx ../jsons/play-0.json -signer user2")
	if !strings.Contains(out, `"Value": "X--------"`) {
		t.Fatal("user2 failed to play")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx ../jsons/play-0.json -signer user1")
	if !strings.Contains(out, `"Value": "that index is taken"`) {
		t.Fatal("user1 was not blocked")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx ../jsons/play-1.json -signer user1")
	if !strings.Contains(out, `"Value": "XO-------"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx ../jsons/play-3.json -signer user2")
	if !strings.Contains(out, `"Value": "XO-X-----"`) {
		t.Fatal("updating board failed")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx ../jsons/play-4.json -signer user1")
	if !strings.Contains(out, `"Value": "XO-XO----"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx ../jsons/play-6.json -signer user2")
	if !strings.Contains(out, `"Value": "congrats, you won"`) {
		t.Fatal("win was not recognized")
	}

	//the game should have ended, lets start a new one
	//add user1 to the pool
	out = gammaCli.Run("send-tx ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}

	//add user2 to the game
	out = gammaCli.Run("send-tx ../jsons/startGame.json -signer user2")
	if !strings.Contains(out, `"Value": "started new game"`) {
		t.Fatal("adding user2 failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx ../jsons/play-0.json -signer user2")
	if !strings.Contains(out, `"Value": "X--------"`) {
		t.Fatal("user2 failed to play")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx ../jsons/play-1.json -signer user1")
	if !strings.Contains(out, `"Value": "XO-------"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx ../jsons/play-4.json -signer user2")
	if !strings.Contains(out, `"Value": "XO--X----"`) {
		t.Fatal("updating board failed")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx ../jsons/play-7.json -signer user1")
	if !strings.Contains(out, `"Value": "XO--X--O-"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx ../jsons/play-8.json -signer user2")
	if !strings.Contains(out, `"Value": "congrats, you won"`) {
		t.Fatal("win was not recognized")
	}

	//the game should have ended, lets start a new one
	//add user1 to the pool
	out = gammaCli.Run("send-tx ../jsons/startGame.json -signer user1")
	if !strings.Contains(out, `"Value": "added to pool"`) {
		t.Fatal("adding user1 failed")
	}

	//add user2 to the game
	out = gammaCli.Run("send-tx ../jsons/startGame.json -signer user2")
	if !strings.Contains(out, `"Value": "started new game"`) {
		t.Fatal("adding user2 failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx ../jsons/play-0.json -signer user2")
	if !strings.Contains(out, `"Value": "X--------"`) {
		t.Fatal("updating board failed")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx ../jsons/play-5.json -signer user1")
	if !strings.Contains(out, `"Value": "X----O---"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx ../jsons/play-1.json -signer user2")
	if !strings.Contains(out, `"Value": "XX---O---"`) {
		t.Fatal("updating board failed")
	}

	//user1 should be stopped from playing, the spot is taken
	out = gammaCli.Run("send-tx ../jsons/play-7.json -signer user1")
	if !strings.Contains(out, `"Value": "XX---O-O-"`) {
		t.Fatal("updating board failed")
	}

	//user2 should play
	out = gammaCli.Run("send-tx ../jsons/play-2.json -signer user2")
	if !strings.Contains(out, `"Value": "congrats, you won"`) {
		t.Fatal("win was not recognized")
	}

}
