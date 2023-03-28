package util

import (
	"context"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestPulumiInstallLogin(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

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
}

func TestPulumiGoInstall(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

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
}

func TestPulumiOnPath(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

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
}

func TestEnvVariable(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

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
}

func TestPulumiInstall(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

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
}
