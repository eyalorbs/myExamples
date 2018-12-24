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
	out := gammaCli.Run("deploy ../token.go -name ERC20Token -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("deploy failed")
	}

	//check if the total amount is okay
	out = gammaCli.Run("run-query ../jsons/totalSupply.json ")
	if !strings.Contains(out, `"Value": "1000000000000000000"`) {
		t.Fatal("total supply failed")
	}

	//check if user1 got all of the tokens
	out = gammaCli.Run("run-query ../jsons/balanceOf-user1.json")
	if !strings.Contains(out, `"Value": "1000000000000000000"`) {
		t.Fatal("initial user1 balance failed")
	}

	//check if approval of user2 is ok
	out = gammaCli.Run("send-tx ../jsons/approve-user2-50.json -signer user1")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("approval of user2 failed")
	}

	//check the allowance of user2 after approval
	out = gammaCli.Run("run-query ../jsons/allowance-user1-user2.json")
	if !strings.Contains(out, `"Value": "50"`) {
		t.Fatal("valued of approval of user2 failed")
	}

	//check the allowance of user 3
	out = gammaCli.Run("run-query ../jsons/allowance-user1-user3.json -signer user3")
	if !strings.Contains(out, `"Value": "0"`) {
		t.Fatal("user3 appears to have an allowance from user1")
	}

	//check if transferFrom was successful
	out = gammaCli.Run("send-tx ../jsons/transferFrom-user1-to-user3.json -signer user2")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("transfer From failed")
	}

	//check the balance of user3 after transfer
	out = gammaCli.Run("run-query ../jsons/balanceOf-user3.json")
	if !strings.Contains(out, `"Value": "20"`) {
		t.Fatal("value of user3 after transferFrom failed")
	}

	//check the balance of user1 after transfer
	out = gammaCli.Run("run-query ../jsons/balanceOf-user1.json")
	if !strings.Contains(out, `"Value": "999999999999999980"`) {
		t.Fatal("value of user1 after transferFrom failed")
	}

	//check if transfer was successful
	out = gammaCli.Run("send-tx ../jsons/transfer-10-user2.json -signer user3")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("transfer to user2 falied")
	}

	//check the balance of user3 after transfer
	out = gammaCli.Run("run-query ../jsons/balanceOf-user3.json")
	if !strings.Contains(out, `"Value": "10"`) {
		t.Fatal("value of user3 after transferFrom failed")
	}

	//check the balance of user2 after transfer
	out = gammaCli.Run("run-query ../jsons/balanceOf-user2.json")
	if !strings.Contains(out, `"Value": "10"`) {
		t.Fatal("value of user1 after transferFrom failed")
	}

}
