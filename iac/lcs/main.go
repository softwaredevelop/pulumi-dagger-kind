//revive:disable:package-comments,exported
package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/debug"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
)

func main() {
	ctx := context.Background()

	stackName := "lcs"
	workDir := filepath.Join("src")

	stack, err := auto.NewStackLocalSource(ctx, stackName, workDir, auto.EnvVars(map[string]string{
		"PULUMI_CONFIG_PASSPHRASE": "",
		"PULUMI_SKIP_UPDATE_CHECK": "true",
	}))
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackName + " already exists")
	}
	if err != nil && !auto.IsCreateStack409Error(err) {
		panic(err)
	}

	ss, err := stack.Workspace().ListStacks(ctx)
	if err != nil {
		panic(err)
	}

	for _, s := range ss {
		log.Println(s.Name)
	}

	prev, err := stack.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(prev.StdOut)
}
