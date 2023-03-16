package main

import (
	"ci/util"
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestPulumiPreviewLocalSource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	c = c.Pipeline("cd")
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
	p := filepath.Join(dir, "..", "lcs")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	id, err = PulumiInstall(c, id).
		Pipeline("pulumi").
		WithWorkdir(mountedDir).
		WithExec([]string{"pulumi", "login", "--local"}).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	out, err := c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-local-source").
		WithWorkdir(mountedDir).
		WithExec([]string{"go", "run", "-v", "main.go"}).
		Stdout(ctx)
	require.NoError(t, err)
	log.Println(out)
}

func TestPulumiGoPreviewInlineSource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	c = c.Pipeline("cd")
	require.NotNil(t, c)

	id, err := c.
		Container().
		From("busybox:glibc").
		// From("golang:alpine").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..", "ils")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	id, err = PulumiGoInstall(c, id).
		Pipeline("pulumi").
		WithWorkdir(mountedDir).
		WithExec([]string{"pulumi", "login", "--local"}).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-inline-source").
		WithWorkdir(mountedDir).
		WithExec([]string{"go", "run", "-v", "main.go"}).
		Stdout(ctx)
	require.Error(t, err)
}

func TestPulumiGoPreviewLocalSourceHost(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	c = c.Pipeline("cd")
	require.NotNil(t, c)

	id, err := c.
		Container().
		From("busybox:glibc").
		// From("golang:alpine").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..", "lsh")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	id, err = PulumiGoInstall(c, id).
		Pipeline("pulumi").
		WithWorkdir(mountedDir).
		WithExec([]string{"pulumi", "login", "--local"}).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-local-source-host").
		WithWorkdir(mountedDir).
		WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
		WithExec([]string{"pulumi", "stack",
			"init",
			"--stack=lsh",
			"--color=always",
			"--non-interactive",
		}).
		WithExec([]string{"pulumi", "stack",
			"ls",
			"--color=always",
			"--non-interactive",
		}).
		WithExec([]string{"pulumi", "preview",
			"--debug",
			"--color=always",
			"--non-interactive",
		}).
		Stdout(ctx)
	require.Error(t, err)
}

func TestPulumiPreviewInlineSource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	c = c.Pipeline("cd")
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
	p := filepath.Join(dir, "..", "ils")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	id, err = PulumiInstall(c, id).
		Pipeline("pulumi").
		WithWorkdir(mountedDir).
		WithExec([]string{"pulumi", "login", "--local"}).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-inline-source").
		WithWorkdir(mountedDir).
		WithExec([]string{"go", "run", "-v", "main.go"}).
		Stdout(ctx)
	require.NoError(t, err)
}

func TestPulumiPreviewLocalSourceHost(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	c = c.Pipeline("cd")
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
	p := filepath.Join(dir, "..", "lsh")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	id, err = PulumiInstall(c, id).
		Pipeline("pulumi").
		WithWorkdir(mountedDir).
		WithExec([]string{"pulumi", "login", "--local"}).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-local-source-host").
		WithWorkdir(mountedDir).
		WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
		WithExec([]string{"pulumi", "stack",
			"init",
			"--stack=lsh",
			"--color=always",
			"--non-interactive",
		}).
		WithExec([]string{"pulumi", "stack",
			"ls",
			"--color=always",
			"--non-interactive",
		}).
		WithExec([]string{"pulumi", "preview",
			"--debug",
			"--color=always",
			"--non-interactive",
		}).
		Stdout(ctx)
	require.NoError(t, err)

	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-local-source-host").
		WithWorkdir(mountedDir).
		WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
		WithExec([]string{"pulumi", "stack",
			"init",
			"--stack=lsh",
			"--color=always",
			"--non-interactive",
		}).
		WithExec([]string{"pulumi", "up",
			"--debug",
			"--skip-preview",
			"--color=always",
			"--non-interactive",
		}).
		Stdout(ctx)
	require.NoError(t, err)
}

func TestMountedHostParentDirectory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

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
}

func TestContainerID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

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
}

func TestErrorMessage(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	_, err = c.
		Container().
		From("fake.invalid").
		ID(ctx)
	require.Error(t, err)
	require.ErrorContains(t, err, "not exist")
}

func TestConnect(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()
}
