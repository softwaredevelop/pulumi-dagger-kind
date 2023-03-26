package lint

import (
	"context"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestEc(t *testing.T) {
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

	container = Ec(c, id)
	require.NotNil(t, container)
	out, err := container.
		WithExec([]string{"ls", "/usr/bin/editorconfig-checker"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/editorconfig-checker\n", out)
}

func TestEcVersion(t *testing.T) {
	t.Parallel()

	version, err := ecVersion()
	require.NoError(t, err)
	require.Contains(t, version, ".")
}

func TestEc2(t *testing.T) {
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

	container = Ec2(c, id)
	require.NotNil(t, container)
	out, err := container.
		WithExec([]string{"ls", "/usr/bin/ec"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/ec\n", out)

	_, err = container.
		WithWorkdir("/tmp").
		WithExec([]string{"ec", "-debug"}).
		Stdout(ctx)
	require.NoError(t, err)
}

func TestEc1(t *testing.T) {
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

	container = Ec1(c, id)
	require.NotNil(t, container)
	out, err := container.
		WithExec([]string{"ls", "/usr/bin/ec"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/ec\n", out)

	_, err = container.
		WithWorkdir("/tmp").
		WithExec([]string{"ec", "-debug"}).
		Stdout(ctx)
	require.NoError(t, err)
}

func TestContainerIDBusybox(t *testing.T) {
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

	out, err := c.Container(dagger.ContainerOpts{ID: id}).
		WithExec([]string{"busybox"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Contains(t, out, "BusyBox")
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

func TestGitCloneFileContent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	repoURL := "https://github.com/editorconfig-checker/editorconfig-checker.git"
	fileName := "VERSION"
	content, err := gitCloneFileContent(ctx, c, repoURL, fileName)
	require.NoError(t, err)
	require.Contains(t, content, ".")
}

func TestGitClone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	repoURL := "https://github.com/editorconfig-checker/editorconfig-checker.git"
	repo := gitClone(c, repoURL)
	require.NotNil(t, repo)
}
