//revive:disable:package-comments,exported
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/debug"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

func main() {
	ctx := context.Background()

	stackName := "lcs-local"
	// workDir := filepath.Join("src")
	workDir := filepath.Join("localproject")

	stack, err := auto.NewStackLocalSource(ctx, stackName, workDir)
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackName + " already exists")
	}
	if err != nil && !auto.IsCreateStack409Error(err) {
		panic(err)
	}

	pat := os.Getenv("PULUMI_ACCESS_TOKEN")
	err = stack.Workspace().SetEnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_CONFIG_PASSPHRASE": "",
		"PULUMI_ACCESS_TOKEN":      pat,
	})
	if err != nil {
		panic(err)
	}

	ss, err := stack.Workspace().ListStacks(ctx)
	if err != nil {
		panic(err)
	}

	contains := false
	for _, s := range ss {
		if s.Name == stackName {
			contains = true
		}
	}
	if !contains {
		panic(stackName + "stack not found")
	}

	ght := os.Getenv("GITHUB:TOKEN")
	gho := os.Getenv("GITHUB:OWNER")
	err = stack.SetAllConfig(ctx, auto.ConfigMap{
		"github:token": auto.ConfigValue{
			Value:  ght,
			Secret: true,
		},
		"github:owner": auto.ConfigValue{
			Value:  gho,
			Secret: true,
		},
	})
	if err != nil {
		panic(err)
	}

	prev, err := stack.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(prev.StdOut)

	up, err := stack.Up(ctx, optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(up.StdOut)
}
