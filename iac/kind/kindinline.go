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
	stackNameA := "kindCommandA"
	stackNameB := "kindCommandB"
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

	stackA, err := auto.NewStackInlineSource(ctx, stackNameA, prj.Name.String(), func(pCtx *pulumi.Context) error {
		pCtx.Export("outputAA", pulumi.String("valueAA"))
		return nil
	}, auto.EnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_CONFIG_PASSPHRASE": "",
	}))
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackNameA + " already exists")
		stackA, err = auto.UpsertStackInlineSource(ctx, stackNameA, prj.Name.String(), func(pCtx *pulumi.Context) error {
			pCtx.Export("outputAA", pulumi.String("valueAA"))
			return nil
		}, auto.EnvVars(map[string]string{
			"PULUMI_SKIP_UPDATE_CHECK": "true",
			"PULUMI_CONFIG_PASSPHRASE": "",
		}))
		if err != nil {
			panic(err)
		}
	}
	if err != nil && !auto.IsCreateStack409Error(err) {
		panic(err)
	}

	stackB, err := auto.NewStackInlineSource(ctx, stackNameB, prj.Name.String(), func(pCtx *pulumi.Context) error {
		pCtx.Export("outputBB", pulumi.String("valueBB"))
		return nil
	}, auto.EnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_CONFIG_PASSPHRASE": "",
	}))
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackNameB + " already exists")
		stackA, err = auto.UpsertStackInlineSource(ctx, stackNameB, prj.Name.String(), func(pCtx *pulumi.Context) error {
			pCtx.Export("outputBB", pulumi.String("valueBB"))
			return nil
		}, auto.EnvVars(map[string]string{
			"PULUMI_SKIP_UPDATE_CHECK": "true",
			"PULUMI_CONFIG_PASSPHRASE": "",
		}))
		if err != nil {
			panic(err)
		}
	}
	if err != nil && !auto.IsCreateStack409Error(err) {
		panic(err)
	}

	if destroy {

		drstA, err := stackA.Destroy(ctx, optdestroy.Message("Successfully destroyed stack : "+stackNameA))
		if err != nil {
			panic(err)
		}
		log.Println(drstA.Summary.Kind + " " + drstA.Summary.Message)

		drstB, err := stackB.Destroy(ctx, optdestroy.Message("Successfully destroyed stack : "+stackNameB))
		if err != nil {
			panic(err)
		}
		log.Println(drstB.Summary.Kind + " " + drstB.Summary.Message)

		return
	}

	stackA.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		kv, err := local.NewCommand(pCtx, "kind-version", &local.CommandArgs{
			Create: pulumi.StringPtrInput(pulumi.String("kind version")),
		})
		if err != nil {
			return err
		}

		clusterName := "testcluster"
		_, err = local.NewCommand(pCtx, "kind-create", &local.CommandArgs{
			Create: pulumi.StringPtrInput(pulumi.String("kind create cluster --name=" + clusterName)),
		})
		if err != nil {
			return err
		}

		pCtx.Export("kindVersion", kv.Stdout)

		return nil

	})

	stackB.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		clusterName := "testcluster"
		_, err = local.NewCommand(pCtx, "kind-delete", &local.CommandArgs{
			Create: pulumi.StringPtrInput(pulumi.String("kind delete cluster --name=" + clusterName)),
		})
		if err != nil {
			return err
		}

		return nil

	})

	prev, err := stackA.Preview(ctx, optpreview.Message("Preview stack "+stackNameA), optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(prev.StdOut)

	refr, err := stackA.Refresh(ctx, optrefresh.Message("Refresh stack "+stackNameA), optrefresh.ProgressStreams(os.Stdout, os.Stderr))
	if err != nil {
		panic(err)
	}
	log.Println(refr.StdOut)

	up, err := stackA.Up(ctx, optup.Message("Update stack "+stackNameA), optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(up.StdOut)

	prev, err = stackB.Preview(ctx, optpreview.Message("Preview stack "+stackNameB), optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(prev.StdOut)

	refr, err = stackB.Refresh(ctx, optrefresh.Message("Refresh stack "+stackNameB), optrefresh.ProgressStreams(os.Stdout, os.Stderr))
	if err != nil {
		panic(err)
	}
	log.Println(refr.StdOut)

	up, err = stackB.Up(ctx, optup.Message("Update stack "+stackNameB), optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(up.StdOut)

}
