package test

import (
	"github.com/orbs-network/orbs-contract-sdk/go/testing/gamma"
	"strconv"
	"strings"
	"testing"
)

func TestDeploy(t *testing.T) {
	gammaCli := gamma.Cli().Start()
	defer gammaCli.Stop()

	for i := 1; i <= 11; i++ {
		message := "deploy ../counter.go -env testnet42 -name MyCounter" + strconv.Itoa(i)
		out := gammaCli.Run(message)
		if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
			t.Fatal("deploy failed")
		}
	}
}
