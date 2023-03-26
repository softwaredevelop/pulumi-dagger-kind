package lint

import (
	"context"
	"os"
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

	container := c.Container().From("busybox:uclibc")
	require.NotNil(t, container)

	id, err := container.ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	container = Revive(c, id)
	require.NotNil(t, container)
	out, err := container.
		WithExec([]string{"ls", "/usr/bin/revive"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/revive\n", out)
}

func TestGoLint(t *testing.T) {
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

	container = GoLint(c, id)
	require.NotNil(t, container)
	out, err := container.
		WithExec([]string{"ls", "/usr/bin/golint"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/golint\n", out)

	out, err = container.
		WithExec([]string{"ls", "/usr/bin/fgt"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/fgt\n", out)
}
