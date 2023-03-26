package lint

import (
	"context"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestAl(t *testing.T) {
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

	container = Al(c, id)
	require.NotNil(t, container)
	out, err := container.
		WithExec([]string{"ls", "/usr/bin/actionlint"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/actionlint\n", out)

	out, err = container.
		WithExec([]string{"ls", "/usr/bin/shellcheck"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/shellcheck\n", out)
}
