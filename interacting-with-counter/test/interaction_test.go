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
	out := gammaCli.Run("deploy -name interact -code ../contract.go")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("deploy failed")
	}

	out = gammaCli.Run("deploy -name MyCounter -code ../../counter/counter.go")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("deploy failed")
	}

	out = gammaCli.Run("read -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "0"`) {
		t.Fatal("read fail")
	}

	out = gammaCli.Run("send-tx -i ../jsons/add-25.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("failed to add")
	}

	out = gammaCli.Run("read -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "25"`) {
		t.Fatal("read fail")
	}

	out = gammaCli.Run("send-tx -i ../jsons/add-25.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("failed to add")
	}

	out = gammaCli.Run("read -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "50"`) {
		t.Fatal("read fail")
	}

	out = gammaCli.Run("send-tx -i ../jsons/add-25.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("failed to add")
	}

	out = gammaCli.Run("read -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "75"`) {
		t.Fatal("read fail")
	}

}
