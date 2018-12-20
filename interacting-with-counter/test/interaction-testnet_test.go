package test

import (
	"github.com/orbs-network/orbs-contract-sdk/go/testing/gamma"
	"strings"
	"testing"
)

func Test_testNet(t *testing.T) {
	gammaCli := gamma.Cli().Start()
	defer gammaCli.Stop()

	//no need to deploy, contract is already deployed

	out := gammaCli.Run("read -env testnet42 -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "0"`) {
		t.Fatal("read fail")
	}

	out = gammaCli.Run("send-tx -env testnet42 -i ../jsons/add-25.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("failed to add")
	}

	out = gammaCli.Run("read -env testnet42 -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "25"`) {
		t.Fatal("read fail")
	}

	out = gammaCli.Run("send-tx -env testnet42 -i ../jsons/add-25.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("failed to add")
	}

	out = gammaCli.Run("read -env testnet42 -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "50"`) {
		t.Fatal("read fail")
	}

	out = gammaCli.Run("send-tx -env testnet42 -i ../jsons/add-25.json")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("failed to add")
	}

	out = gammaCli.Run("read -env testnet42 -i ../jsons/get.json")
	if !strings.Contains(out, `"Value": "75"`) {
		t.Fatal("read fail")
	}

}
