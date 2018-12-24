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
	out := gammaCli.Run("deploy ../helloWorld.go -env testnet42 -signer user1 -name helloWorld12")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("deploy failed")
	}

	//check output
	out = gammaCli.Run("run-query ../jsons/greet.json -env testnet42")
	if !strings.Contains(out, `"Value": "hello world!"`) {
		t.Fatal("greeting failed")
	}
}
