package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/frezbo/pulumi-provider-kind/sdk/v3/go/kind/cluster"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestKind(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	projectName := "testproject"
	t.Run("Test_inline_source_kind_command", func(t *testing.T) {
		t.Parallel()

		stackName := "testInlineSourceKind"
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
		}, auto.EnvVars(map[string]string{
			"PULUMI_SKIP_UPDATE_CHECK": "true",
			"PULUMI_CONFIG_PASSPHRASE": "",
		}))

		defer func() {
			err = s.Workspace().RemoveStack(ctx, s.Name())
			require.NoError(t, err)
		}()

		s.Workspace().SetProgram(func(pCtx *pulumi.Context) error {
			clusterName := "testclusterkind"
			cluster, err := cluster.NewCluster(pCtx, clusterName, &cluster.ClusterArgs{
				Name: pulumi.String(clusterName),
			})
			require.NoError(t, err)
			require.NotNil(t, cluster)

			return nil
		})
	})
	t.Run("Test_new_stack_inline_source_destroy", func(t *testing.T) {
		t.Parallel()

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
		}, auto.EnvVars(map[string]string{
			"PULUMI_SKIP_UPDATE_CHECK": "true",
			"PULUMI_CONFIG_PASSPHRASE": "",
		}))

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
	})
	t.Run("Test_new_stack_inline_source_project", func(t *testing.T) {
		t.Parallel()

		stackName := "testInlineSourceProject"
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
		}, auto.EnvVars(map[string]string{
			"PULUMI_SKIP_UPDATE_CHECK": "true",
			"PULUMI_CONFIG_PASSPHRASE": "",
		}))
		require.NoError(t, err)
		require.NotNil(t, s)

		defer func() {
			err = s.Workspace().RemoveStack(ctx, s.Name())
			require.NoError(t, err)
		}()

		prj, err = s.Workspace().ProjectSettings(ctx)
		require.NoError(t, err)
		require.Contains(t, prj.Name.String(), "testproject")
		require.Contains(t, prj.Runtime.Name(), "go")

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

		values, err := s.GetAllConfig(ctx)
		require.NoError(t, err)

		require.Equal(t, "def", values["bar:token"].Value)
		require.Equal(t, "xyz", values["buzz:owner"].Value)
	})
	t.Run("Test_new_stack_inline_source_secrets", func(t *testing.T) {
		t.Parallel()

		stackName := "testInlineSourceSecrets"
		s, err := auto.NewStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
			return nil
		}, auto.EnvVars(map[string]string{
			"PULUMI_SKIP_UPDATE_CHECK": "true",
			"PULUMI_CONFIG_PASSPHRASE": "",
		}))

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
	})
	t.Run("Test_new_stack_inline_source_config", func(t *testing.T) {
		t.Parallel()

		stackName := "testInlineSourceConfig"
		s, err := auto.NewStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
			return nil
		}, auto.EnvVars(map[string]string{
			"PULUMI_SKIP_UPDATE_CHECK": "true",
			"PULUMI_CONFIG_PASSPHRASE": "",
		}))
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
	})
	t.Run("Test_new_stack_inline_source", func(t *testing.T) {
		t.Parallel()
		stackName := "testInlineSourceWorkspaceEnvVars"
		projectName := "testproject"

		s, err := auto.NewStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
			return nil
		}, auto.EnvVars(map[string]string{
			"PULUMI_SKIP_UPDATE_CHECK": "true",
			"PULUMI_CONFIG_PASSPHRASE": "",
		}))
		require.NoError(t, err)
		require.NotNil(t, s)

		defer func() {
			err = s.Workspace().RemoveStack(ctx, s.Name())
			require.NoError(t, err)
		}()

		require.Equal(t, stackName, s.Name())

		ss, err := s.Workspace().ListStacks(ctx)
		require.NoError(t, err)
		require.NotNil(t, ss)
		sss := []string{}
		for _, s := range ss {
			sss = append(sss, s.Name)
		}
		require.Contains(t, sss, stackName)

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
	})
}
