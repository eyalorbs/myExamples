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
	out := gammaCli.Run("deploy ../counter.go -name MyCounter12")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("deploy failed")
	}

	//check output
	out = gammaCli.Run("run-query ../jsons/get.json")
	if !strings.Contains(out, `"Value": "0"`) {
		t.Fatal("initial count failed")
	}

	//add to the count and check
	out = gammaCli.Run("send-tx  ../jsons/add-25.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal(`"ExecutionResult": "SUCCESS"`)
	}

	//check output
	out = gammaCli.Run("run-query ../jsons/get.json")
	if !strings.Contains(out, `"Value": "25"`) {
		t.Fatal("initial count failed")
	}

	//add to the count and check
	out = gammaCli.Run("send-tx ../jsons/add-10.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal(`"ExecutionResult": "SUCCESS"`)
	}

	//check output
	out = gammaCli.Run("run-query ../jsons/get.json")
	if !strings.Contains(out, `"Value": "35"`) {
		t.Fatal("initial count failed")
	}
}
