package util

import (
	"context"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestPulumi(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	t.Run("Test_pulumi_scratch_install", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		c = c.Pipeline("test-pulumi-scratch-install")
		id, err := c.Container().ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		container := c.Container(dagger.ContainerOpts{ID: id})
		require.NotNil(t, container)

		container = P(c, id)
		require.NotNil(t, container)
		out, err := container.
			WithExec([]string{"/pulumi", "version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "v")
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

		container = PulumiInstall(c, id)
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

		container = PulumiGoInstall(c, id)
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
	t.Run("Test_pulumi_on_path", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()

		c = c.Pipeline("test-pulumi-on-path")
		container := c.Container().From("busybox:uclibc")
		require.NotNil(t, container)

		id, err := container.ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		container = PulumiInstall(c, id)
		require.NotNil(t, container)

		_, err = container.
			WithExec([]string{"pulumi", "login", "--local"}).
			Stdout(ctx)
		require.NoError(t, err)
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
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
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

		container = PulumiInstall(c, id)
		require.NotNil(t, container)
		out, err := container.
			WithExec([]string{"ls", "/usr/bin/pulumi"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Equal(t, "/usr/bin/pulumi\n", out)
	})
}
