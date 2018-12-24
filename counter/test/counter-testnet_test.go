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
	out := gammaCli.Run("run-query ../jsons/get.json -env testnet42")
	if !strings.Contains(out, `"Value": "135"`) {
		t.Fatal("initial count failed")
	}

	//add to the count and check
	out = gammaCli.Run("send-tx ../jsons/add-10.json -env testnet42")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal(`"ExecutionResult": "SUCCESS"`)
	}

	//check output
	out = gammaCli.Run("run-query ../jsons/get.json -env testnet42")
	if !strings.Contains(out, `"Value": "145"`) {
		t.Fatal("initial count failed")
	}

	//add to the count and check
	out = gammaCli.Run("send-tx ../jsons/add-25.json -env testnet42")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal(`"ExecutionResult": "SUCCESS"`)
	}

	//check output
	out = gammaCli.Run("run-query ../jsons/get.json -env testnet42")
	if !strings.Contains(out, `"Value": "170"`) {
		t.Fatal("initial count failed")
	}
}
