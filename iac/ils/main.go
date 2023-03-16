//revive:disable:package-comments,exported
package main

import (
	"context"
	"log"

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
	desc := "A minimal inline source Go Pulumi program"
	_, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:        "iac",
		Runtime:     workspace.NewProjectRuntimeInfo("go", nil),
		Description: &desc,
	}))
	if err != nil {
		panic(err)
	}

	stack, _ := auto.NewStackInlineSource(ctx, stackName, "iac", func(pCtx *pulumi.Context) error {
		return nil
	}, auto.EnvVars(map[string]string{
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
