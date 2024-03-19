package main

import (
	"context"
	"log"
	"testing"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/debug"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optrefresh"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optremove"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestUpsertStackInlineSourceRefresh(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	projectName := "testproject"
	stackNameA := "testUpsertStackInlineSourceRefreshA"
	desc := "A inline source Go Pulumi program Test"
	ws, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:        tokens.PackageName(projectName),
		Runtime:     workspace.NewProjectRuntimeInfo("go", nil),
		Description: &desc,
	}))
	require.NoError(t, err)
	require.NotNil(t, ws)

	prj, err := ws.ProjectSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, prj)

	s, err := auto.UpsertStackInlineSource(ctx, stackNameA, projectName, func(pCtx *pulumi.Context) error {
		pCtx.Export("outputAA", pulumi.String("valueAA"))
		return nil
	}, auto.EnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_CONFIG_PASSPHRASE": "",
	}))

	require.NoError(t, err)
	require.NotNil(t, s)

	defer func() {
		dr, err := s.Destroy(ctx, optdestroy.Message("Successfully destroyed stack :"+stackNameA))
		require.NoError(t, err)
		log.Println(dr.Summary.Kind + " " + dr.Summary.Message)
		err = s.Workspace().RemoveStack(ctx, s.Name(), optremove.Force())
		require.NoError(t, err)
	}()

	err = s.SetAllConfig(ctx, auto.ConfigMap{
		"bar:token": auto.ConfigValue{
			Value:  "def",
			Secret: true,
		},
		"buzz:owner": auto.ConfigValue{
			Value:  "xyz",
			Secret: true,
		},
	})
	require.NoError(t, err)

	rr, err := s.Refresh(ctx, optrefresh.Message("Refresh stack "+stackNameA))
	require.NoError(t, err)
	require.NotNil(t, rr)
	log.Println(rr.Summary.Kind + " " + rr.Summary.Message)

	values, err := s.GetAllConfig(ctx)
	require.NoError(t, err)

	require.Equal(t, "def", values["bar:token"].Value)
	require.Equal(t, "xyz", values["buzz:owner"].Value)

	s.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		hello, err := local.NewCommand(pCtx, "hello", &local.CommandArgs{
			Create: pulumi.String("echo \"Hello Pulumi\""),
		})
		if err != nil {
			return err
		}

		pCtx.Export("hello", hello.Stdout)

		return nil
	})

	prev, err := s.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(prev.StdOut)

	up, err := s.Up(ctx, optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(up.StdOut)

	stackNameB := "testUpsertStackInlineSourceRefreshB"

	require.NoError(t, err)
	require.NotNil(t, ws)

	prj, err = ws.ProjectSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, prj)

	ss, err := auto.UpsertStackInlineSource(ctx, stackNameB, projectName, func(pCtx *pulumi.Context) error {
		pCtx.Export("outputBB", pulumi.String("valueBB"))
		return nil
	}, auto.EnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_CONFIG_PASSPHRASE": "",
	}))
	require.NoError(t, err)
	require.NotNil(t, ss)

	defer func() {
		dr, err := ss.Destroy(ctx, optdestroy.Message("Successfully destroyed stack :"+stackNameB))
		require.NoError(t, err)
		log.Println(dr.Summary.Kind + " " + dr.Summary.Message)
		err = ss.Workspace().RemoveStack(ctx, ss.Name(), optremove.Force())
		require.NoError(t, err)
	}()

	ss.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		stackReff, err := pulumi.NewStackReference(pCtx, stackNameA, nil)
		require.NoError(t, err)
		require.NotNil(t, stackReff)

		outputValue := stackReff.GetOutput(pulumi.String("hello"))
		require.NotNil(t, outputValue)
		log.Println(outputValue)

		return nil

	})

	prev, err = ss.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(prev.StdOut)

	up, err = ss.Up(ctx, optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(up.StdOut)
}
