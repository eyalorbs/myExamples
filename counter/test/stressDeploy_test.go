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

	for i := 117; i <= 5000; i++ {
		message := "deploy -env testnet42 -name MyCounter" + strconv.Itoa(i) + " -code ../counter.go"
		out := gammaCli.Run(message)
		if !strings.Contains(out, `"ExecutionResult": "SUCCESS"`) {
			t.Fatal("deploy failed")
		}
	}
}
