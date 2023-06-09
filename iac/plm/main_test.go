package main

import (
	"context"
	"os"
	"path/filepath"
	"plm/util"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestPulumiInlineSource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	t.Run("Test_pulumi_inline_source_kind", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		c = c.Pipeline("cd-inline-source-kind-test")
		require.NotNil(t, c)

		id, err := c.
			Container().
			// From("busybox:glibc").
			From("golang:alpine").
			WithMountedTemp("/mountedtmp").
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		dir, _ := os.Getwd()
		p := filepath.Join(dir, "..", "kind")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		mountedDir := "/mountedtmp"
		id, err = util.MountedHostDirectory(c, id, p, mountedDir).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		id, err = util.PulumiInstall(c, id).
			Pipeline("pulumi").
			WithWorkdir(mountedDir).
			WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
			WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
			WithExec([]string{"pulumi", "login", "--local"}).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		// reMatching := "TestInlineSourceKindCommand$"
		// _, err = c.Container(dagger.ContainerOpts{ID: id}).
		// 	Pipeline("pulumi-inline-source-kind-test1").
		// 	WithWorkdir(mountedDir).
		// 	WithExec([]string{"go", "test", "-v", "-run", reMatching}).
		// 	Stdout(ctx)
		// require.Error(t, err)
		// require.Contains(t, err.Error(), "exit status 1")
	})
	t.Run("Test_mounted_host_parent_directory", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		container := c.Container().From("busybox:glibc")
		require.NotNil(t, container)

		dir, _ := os.Getwd()
		p := filepath.Join(dir, "..")
		require.NoError(t, err)
		require.NotEmpty(t, p)
		id, err := container.
			WithMountedDirectory("/mountedtmp", c.Host().Directory(p)).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		mntdir, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"ls", "/mountedtmp"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, mntdir)

		lsdir, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithWorkdir("/mountedtmp").
			WithExec([]string{"ls", "-la"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, lsdir)
		require.Contains(t, lsdir, "kind")
	})
	t.Run("Test_container_ID", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		container := c.Container().From("alpine")
		require.NotNil(t, container)

		id, err := container.ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		releaseName, err := c.Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"/bin/sh", "-c", "cat /etc/os-release | awk -F= '/^NAME/ {print $2}' | tr -d '\"'"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Equal(t, "Alpine Linux\n", releaseName)
	})
	t.Run("Test_error_message", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		_, err = c.
			Container().
			From("fake.invalid").
			ID(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "not exist")
	})
	t.Run("Test_pulumi_install_login", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		c = c.Pipeline("test-pulumi-install-login")
		container := c.Container().From("busybox:uclibc")
		require.NotNil(t, container)

		id, err := container.ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		container = c.Container(dagger.ContainerOpts{ID: id})
		require.NotNil(t, container)

		container = util.PulumiInstall(c, id)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"ls", "/usr/bin/pulumi"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Equal(t, "/usr/bin/pulumi\n", out)

		out, err = container.
			WithExec([]string{"pulumi", "login", "--local"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, out)
	})
	t.Run("Test_pulumi_go_install", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		c = c.Pipeline("test-pulumi-go-install")
		container := c.Container().From("busybox:glibc")
		require.NotNil(t, container)

		id, err := container.ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		container = util.PulumiGoInstall(c, id)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"go", "version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "version")

		out, err = container.
			WithExec([]string{"pulumi", "version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "v")
	})
	t.Run("Test_env_variable", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		c = c.Pipeline("test-env-variable")
		container := c.Container().From("busybox:glibc")
		require.NotNil(t, container)

		id, err := container.ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		container = c.Container(dagger.ContainerOpts{ID: id})
		require.NotNil(t, container)

		env, err := container.EnvVariable(ctx, "PATH")
		require.NoError(t, err)
		require.Contains(t, env, "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")

		env, err = container.
			WithEnvVariable("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin").
			EnvVariable(ctx, "PATH")
		require.NoError(t, err)
		require.Contains(t, env, "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	})
	t.Run("Test_pulumi_install", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		c = c.Pipeline("test-pulumi-install")
		container := c.Container().From("busybox:uclibc")
		require.NotNil(t, container)

		id, err := container.ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		container = c.Container(dagger.ContainerOpts{ID: id})
		require.NotNil(t, container)

		container = util.PulumiInstall(c, id)
		require.NotNil(t, container)
		out, err := container.
			WithExec([]string{"ls", "/usr/bin/pulumi"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Equal(t, "/usr/bin/pulumi\n", out)
	})
}
