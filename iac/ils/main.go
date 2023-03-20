//revive:disable:package-comments,exported
package main

import (
	"context"
	"log"
	"os"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/debug"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	ctx := context.Background()

	stackName := "ils"
	projectName := "iac"
	desc := "A inline source Go Pulumi program"
	_, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:        "iac",
		Runtime:     workspace.NewProjectRuntimeInfo("go", nil),
		Description: &desc,
	}))
	if err != nil {
		panic(err)
	}

	stack, err := auto.NewStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		return nil
	})
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackName + " already exists")
	}
	if err != nil && !auto.IsCreateStack409Error(err) {
		panic(err)
	}

	err = stack.Workspace().SetEnvVars(map[string]string{
		"PULUMI_CONFIG_PASSPHRASE": "",
		"PULUMI_SKIP_UPDATE_CHECK": "true",
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

	stack.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		hello, err := local.NewCommand(pCtx, "hello", &local.CommandArgs{
			Create: pulumi.String("echo \"Hello Pulumi\""),
		})
		if err != nil {
			return err
		}

		pCtx.Export("hello", hello.Stdout)

		return nil
	})

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
