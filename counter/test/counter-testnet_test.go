package test

import (
	"github.com/orbs-network/orbs-contract-sdk/go/testing/gamma"
	"strings"
	"testing"
)

func Test_testNet(t *testing.T) {
	gammaCli := gamma.Cli().Start()
	defer gammaCli.Stop()

	//no need to deploy contract, it's already deployed

	//check output
	out := gammaCli.Run("read -env testnet42 -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "0"`) {
		t.Fatal("initial count failed")
	}

	//add to the count and check
	out = gammaCli.Run("send-tx -env testnet42 -i ../jsons/add-10.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal(`"ExecutionResult": "SUCCESS"`)
	}

	//check output
	out = gammaCli.Run("read -env testnet42 -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "10"`) {
		t.Fatal("initial count failed")
	}

	//add to the count and check
	out = gammaCli.Run("send-tx -env testnet42 -i ../jsons/add-25.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal(`"ExecutionResult": "SUCCESS"`)
	}

	//check output
	out = gammaCli.Run("read -env testnet42 -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "35"`) {
		t.Fatal("initial count failed")
	}
}
