package main

import (
	"ci/lint"
	"ci/util"
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestRevive(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	c = c.Pipeline("ci-test")
	require.NotNil(t, c)

	id, err := c.
		Container().
		From("busybox:uclibc").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	out, err := lint.Revive(c, id).
		Pipeline("revive-test").
		WithWorkdir(mountedDir).
		WithExec([]string{"ls", "-la"}).
		WithExec([]string{"revive", "-set_exit_status", "./..."}).
		Stdout(ctx)
	require.NoError(t, err)
	log.Println(out)
}

func TestClientPipeline(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	p1 := c.Pipeline("lint-test")
	require.NotNil(t, p1)

	id, err := p1.
		Container().
		From("busybox:uclibc").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(p1, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	ec, err := lint.Ec(p1, id).
		Pipeline("ec-test").
		WithWorkdir(mountedDir).
		WithExec([]string{"editorconfig-checker", "-verbose"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, ec)

	al, err := lint.Al(p1, id).
		Pipeline("actionlint-test").
		WithWorkdir(mountedDir).
		WithExec([]string{"actionlint", "-pyflakes=", "-verbose", ".github/workflows/ci.yml"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Empty(t, al)
}

func TestMountedHostRootDirectory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	container := c.Container().From("busybox:uclibc")
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
