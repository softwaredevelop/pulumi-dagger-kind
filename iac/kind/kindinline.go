//revive:disable:package-comments,exported
package main

import (
	"context"
	"log"
	"os"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/debug"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optrefresh"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	destroy := false
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "destroy" {
			destroy = true
		}
	}

	ctx := context.Background()

	projectName := "thesis"
	stackName := "kindCommand"
	desc := "A inline source Go Pulumi program for kind in thesis project."
	ws, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:        tokens.PackageName(projectName),
		Runtime:     workspace.NewProjectRuntimeInfo("go", nil),
		Description: &desc,
	}))
	if err != nil {
		panic(err)
	}

	prj, err := ws.ProjectSettings(ctx)
	if err != nil {
		panic(err)
	}

	stack, err := auto.NewStackInlineSource(ctx, stackName, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackName + " already exists")
		stack, err = auto.UpsertStackInlineSource(ctx, stackName, prj.Name.String(), func(ctx *pulumi.Context) error {
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
	if err != nil && !auto.IsCreateStack409Error(err) {
		panic(err)
	}

	err = stack.Workspace().SetEnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_CONFIG_PASSPHRASE": "",
	})
	if err != nil {
		panic(err)
	}

	if destroy {

		drst, err := stack.Destroy(ctx, optdestroy.Message("Successfully destroyed stack : "+stackName))
		if err != nil {
			panic(err)
		}
		log.Println(drst.Summary.Kind + " " + drst.Summary.Message)

		return
	}

	stack.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		kv, err := local.NewCommand(pCtx, "kind-version", &local.CommandArgs{
			Create: pulumi.StringPtrInput(pulumi.String("kind version")),
		})
		if err != nil {
			return err
		}

		clusterName := "testcluster"
		kc, err := local.NewCommand(pCtx, "kind-create", &local.CommandArgs{
			Create: pulumi.StringPtrInput(pulumi.String("kind create cluster --name=" + clusterName)),
		})
		if err != nil {
			return err
		}

		_, err = local.NewCommand(pCtx, "kind-delete", &local.CommandArgs{
			Create: pulumi.StringPtrInput(pulumi.String("kind delete cluster --name=" + clusterName)),
		})
		if err != nil {
			return err
		}

		pCtx.Export("kindVersion", kv.Stdout)
		pCtx.Export("kindClusterName", kc.Stdout)

		return nil

	})

	prev, err := stack.Preview(ctx, optpreview.Message("Preview stack "+stackName), optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(prev.StdOut)

	refr, err := stack.Refresh(ctx, optrefresh.Message("Refresh stack "+stackName), optrefresh.ProgressStreams(os.Stdout, os.Stderr))
	if err != nil {
		panic(err)
	}
	log.Println(refr.StdOut)

	up, err := stack.Up(ctx, optup.Message("Update stack "+stackName), optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(up.StdOut)

}
