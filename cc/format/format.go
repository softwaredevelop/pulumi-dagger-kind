//revive:disable:package-comments,exported
package format

import (
	"context"

	"dagger.io/dagger"
)

func Gofumpt(c *dagger.Client, id dagger.ContainerID, mountedDir string) error {
	ctx := context.Background()
	_, err := GoFormat(c, id).
		WithWorkdir(mountedDir).
		WithExec([]string{"/gofumpt", "-l", "-w", "."}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GoImports(c *dagger.Client, id dagger.ContainerID, mountedDir string) error {
	ctx := context.Background()
	_, err := GoFormat(c, id).
		WithWorkdir(mountedDir).
		WithExec([]string{"/goimports", "-l", "-w", "."}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GoFormat(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return goFormat(c, id)
}

func goFormat(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("golang:alpine").
		WithExec([]string{
			"go", "install",
			"mvdan.cc/gofumpt@latest",
		}).
		WithExec([]string{
			"go", "install",
			"golang.org/x/tools/cmd/goimports@latest",
		})
	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", install.File("/go/bin/gofumpt")).
		WithFile("/", install.File("/go/bin/goimports"))
}
