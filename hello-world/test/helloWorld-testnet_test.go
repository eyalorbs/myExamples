package test

import (
	"github.com/orbs-network/orbs-contract-sdk/go/testing/gamma"
	"strings"
	"testing"
)

func Test_testNet(t *testing.T) {
	gammaCli := gamma.Cli().Start()
	defer gammaCli.Stop()

	//check output
	out := gammaCli.Run("read -env testnet42 -i ../jsons/greet.json")
	if !strings.Contains(out, `"Value": "hello world!"`) {
		t.Fatal("greeting failed")
	}
}
