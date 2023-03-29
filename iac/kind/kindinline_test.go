package main

import (
	"context"
	"log"
	"os"
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

func TestInlineSourceKindCommand(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	projectName := "kindtestproject"
	stackName := "testKindCommand"
	desc := "A inline source Go Pulumi program for test kind"
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

	s, err := auto.NewStackInlineSource(ctx, stackName, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackName + " already exists")
		s, err = auto.UpsertStackInlineSource(ctx, stackName, prj.Name.String(), func(ctx *pulumi.Context) error {
			return nil
		})
	}
	require.NoError(t, err)
	require.NotNil(t, s)

	defer func() {
		dr, err := s.Destroy(ctx, optdestroy.Message("Successfully destroyed stack :"+stackName))
		require.NoError(t, err)
		log.Println(dr.Summary.Kind + " " + dr.Summary.Message)
		err = s.Workspace().RemoveStack(ctx, s.Name(), optremove.Force())
		require.NoError(t, err)
	}()

	err = s.Workspace().SetEnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_CONFIG_PASSPHRASE": "",
	})

	s.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		kv, err := local.NewCommand(pCtx, "kind-version", &local.CommandArgs{
			Create: pulumi.StringPtrInput(pulumi.String("kind version")),
		})
		require.NoError(t, err)
		require.NotNil(t, kv)

		clusterName := "testcluster"
		kc, err := local.NewCommand(pCtx, "kind-create", &local.CommandArgs{
			Create: pulumi.StringPtrInput(pulumi.String("kind create cluster --name=" + clusterName)),
		})
		require.NoError(t, err)
		require.NotNil(t, kc)

		_, err = local.NewCommand(pCtx, "kind-delete", &local.CommandArgs{
			Create: pulumi.StringPtrInput(pulumi.String("kind delete cluster --name=" + clusterName)),
		})
		require.NoError(t, err)

		pCtx.Export("kindVersion", kv.Stdout)
		pCtx.Export("kindClusterName", kc.Stdout)

		return nil
	})
	require.NoError(t, err)

	prev, err := s.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(prev.StdOut)

	refr, err := s.Refresh(ctx, optrefresh.Message("Refresh stack "+stackName), optrefresh.ProgressStreams(os.Stdout, os.Stderr))
	if err != nil {
		panic(err)
	}
	log.Println(refr.StdOut)

	_, err = s.Up(ctx, optup.Message("Update stack "+stackName), optup.Option(
		optup.DebugLogging(debug.LoggingOptions{
			Debug: true,
		}),
	))
	require.NoError(t, err)

}

func TestNewStackInlineSourceDestroy(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	projectName := "testproject"
	stackName := "testInlineSourceDestroy"

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

	s, err := auto.NewStackInlineSource(ctx, stackName, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})
	require.NoError(t, err)
	require.NotNil(t, s)

	defer func() {
		err = s.Workspace().RemoveStack(ctx, s.Name())
		require.NoError(t, err)
	}()

	os.Args = []string{"test", "destroy"}
	args := os.Args
	if len(args) > 0 {
		if args[0] == "destroy" {
			drst, err := s.Destroy(ctx, optdestroy.Message("Successfully destroyed stack test1"))
			require.NoError(t, err)
			require.NotNil(t, drst)
			log.Println(drst.Summary.Kind + " " + drst.Summary.Message)
		}
	}

	os.Args = []string{"destroy", "test"}
	args = os.Args
	if len(args) > 0 {
		if args[0] == "destroy" {
			drst, err := s.Destroy(ctx, optdestroy.Message("Successfully destroyed stack test2"))
			require.NoError(t, err)
			require.NotNil(t, drst)
			log.Println(drst.Summary.Kind + " " + drst.Summary.Message)
		}
	}

}

func TestNewStackInlineSourceProject(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pon := os.Getenv("PULUMI_ORG_NAME")
	orgName := pon
	projectName := "testproject"
	stackName := "testInlineSourceProject"

	desc := "A inline source Go Pulumi program Test"
	ws, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:        tokens.PackageName(projectName),
		Runtime:     workspace.NewProjectRuntimeInfo("go", nil),
		Description: &desc,
		Config: map[string]workspace.ProjectConfigType{
			"bar:token": {
				Value: "abc",
			},
		},
	}))
	require.NoError(t, err)
	require.NotNil(t, ws)

	prj, err := ws.ProjectSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, prj)

	qualifiedStackName := auto.FullyQualifiedStackName(orgName, prj.Name.String(), stackName)
	require.NotNil(t, qualifiedStackName)
	require.Equal(t, orgName+"/"+projectName+"/"+stackName, qualifiedStackName)

	s, err := auto.NewStackInlineSource(ctx, stackName, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})
	require.NoError(t, err)
	require.NotNil(t, s)

	defer func() {
		err = s.Workspace().RemoveStack(ctx, s.Name())
		require.NoError(t, err)
	}()

	prj, err = s.Workspace().ProjectSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, prj)
	log.Println("project name: " + prj.Name.String())
	log.Println("project runtime: " + prj.Runtime.Name())

	values, err := s.GetAllConfig(ctx)
	require.NoError(t, err)

	for _, s := range values {
		log.Println("config: " + s.Value)
	}

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

	values, err = s.GetAllConfig(ctx)
	require.NoError(t, err)

	require.Equal(t, "def", values["bar:token"].Value)
	require.Equal(t, "xyz", values["buzz:owner"].Value)
}

func TestNewStackInlineSourceSecrets(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	stackName := "testInlineSourceSecrets"
	projectName := "testproject"

	s, err := auto.NewStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		return nil
	}, auto.SecretsProvider("passphrase"))
	require.NoError(t, err)
	require.NotNil(t, s)

	defer func() {
		err = s.Workspace().RemoveStack(ctx, s.Name())
		require.NoError(t, err)
	}()

	f := os.Getenv("FOO")
	b := os.Getenv("BAZ")
	s.Workspace().SetEnvVars(map[string]string{
		"FOO": f,
		"BAZ": b,
	})

	envvars := s.Workspace().GetEnvVars()
	require.Equal(t, f, envvars["FOO"])
	require.Equal(t, b, envvars["BAZ"])

	b = os.Getenv("BAR:TOKEN")
	err = s.SetAllConfig(ctx, auto.ConfigMap{
		"bar:token": auto.ConfigValue{
			Value:  b,
			Secret: true,
		},
		"buzz:owner": auto.ConfigValue{
			Value:  "xyz",
			Secret: true,
		},
	})
	require.NoError(t, err)

	values, err := s.GetAllConfig(ctx)
	require.NoError(t, err)

	require.Equal(t, b, values["bar:token"].Value)
	require.True(t, values["bar:token"].Secret)

	require.Equal(t, "xyz", values["buzz:owner"].Value)
	require.True(t, values["buzz:owner"].Secret)
}

func TestNewStackInlineSourceConfig(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	stackName := "testInlineSourceConfig"
	projectName := "testproject"

	s, err := auto.NewStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		return nil
	})
	require.NoError(t, err)
	require.NotNil(t, s)

	defer func() {
		err = s.Workspace().RemoveStack(ctx, s.Name())
		require.NoError(t, err)
	}()

	require.Equal(t, stackName, s.Name())

	err = s.SetAllConfig(ctx, auto.ConfigMap{
		"bar:token": auto.ConfigValue{
			Value:  "abc",
			Secret: true,
		},
		"buzz:owner": auto.ConfigValue{
			Value:  "xyz",
			Secret: true,
		},
	})
	require.NoError(t, err)

	values, err := s.GetAllConfig(ctx)
	require.NoError(t, err)

	require.Equal(t, "abc", values["bar:token"].Value)
	require.True(t, values["bar:token"].Secret)

	require.Equal(t, "xyz", values["buzz:owner"].Value)
	require.True(t, values["buzz:owner"].Secret)
}

func TestNewStackInlineSourceWorkspaceEnvVars(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	stackName := "testInlineSourceWorkspaceEnvVars"
	projectName := "testproject"

	s, err := auto.NewStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		return nil
	})
	require.NoError(t, err)
	require.NotNil(t, s)

	defer func() {
		err = s.Workspace().RemoveStack(ctx, s.Name())
		require.NoError(t, err)
	}()

	require.Equal(t, stackName, s.Name())

	// ss, err := s.Workspace().ListStacks(ctx)
	// require.NoError(t, err)
	// require.NotNil(t, ss)
	// sss := []string{}
	// for _, s := range ss {
	// 	sss = append(sss, s.Name)
	// }
	// require.Contains(t, sss, stackName)

	err = s.Workspace().SetEnvVars(map[string]string{
		"FOO": "BAR",
		"BAZ": "QUX",
	})
	require.NoError(t, err)

	envvars := s.Workspace().GetEnvVars()
	require.Equal(t, "BAR", envvars["FOO"])
	require.Equal(t, "QUX", envvars["BAZ"])

	s.Workspace().UnsetEnvVar("FOO")
	s.Workspace().UnsetEnvVar("BAZ")

	envvars = s.Workspace().GetEnvVars()
	require.NotContains(t, envvars, "FOO")
	require.NotContains(t, envvars, "BAZ")
}
