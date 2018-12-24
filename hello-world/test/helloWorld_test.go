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
	out := gammaCli.Run("deploy ../helloWorld.go -signer user1 -name helloWorld12")
	if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
		t.Fatal("deploy failed")
	}

	//check output
	out = gammaCli.Run("run-query ../jsons/greet.json")
	if !strings.Contains(out, `"Value": "hello world!"`) {
		t.Fatal("greeting failed")
	}
}
