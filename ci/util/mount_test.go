package util

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestMountedHostDirectoryExcludeFiles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	id, err := c.
		Container().
		From("busybox:uclibc").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..", "..")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	mntdir, err := c.
		Container(dagger.ContainerOpts{ID: id}).
		WithExec([]string{"ls", "-la", "/mountedtmp"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, mntdir)
	log.Println(mntdir)
}

func TestHostDirectoryPath(t *testing.T) {
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
}

func TestMountedTempDirectory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	id, err := c.
		Container().
		From("busybox:uclibc").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)

	mntemp, err := c.
		Container(dagger.ContainerOpts{ID: id}).
		WithExec([]string{"ls", "/"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Contains(t, mntemp, "mountedtmp")

	mntdir, err := c.
		Container(dagger.ContainerOpts{ID: id}).
		WithMountedDirectory("/mountedtmp", c.Host().Directory(".")).
		WithExec([]string{"ls", "/mountedtmp"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, mntdir)
}

func TestMountedHostDirectory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	id, err := c.
		Container().
		From("busybox:uclibc").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)

	container := MountedHostDirectory(c, id, ".", "/mountedtmp")
	require.NotNil(t, container)

	mntdir, err := c.
		Container(dagger.ContainerOpts{ID: id}).
		WithMountedDirectory("/mountedtmp", c.Host().Directory(".")).
		WithExec([]string{"ls", "-la", "/mountedtmp"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, mntdir)
}
