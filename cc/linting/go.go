//revive:disable:package-comments,exported
package linting

import (
	"context"

	"dagger.io/dagger"
)

func Revive(c *dagger.Client, id dagger.ContainerID, mountedDir string) error {
	ctx := context.Background()
	_, err := ReviveL(c, id).
		WithWorkdir(mountedDir).
		WithExec([]string{"/revive", "-set_exit_status", "./..."}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}

func ReviveL(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return revivel(c, id)
}

func revivel(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("golang:alpine").
		WithExec([]string{
			"go", "install",
			"github.com/mgechev/revive@latest",
		})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", install.File("/go/bin/revive"))
}

func GoLint(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return goLint(c, id)
}

func goLint(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("golang:alpine").
		WithExec([]string{
			"go", "install",
			"golang.org/x/lint/golint@latest",
		}).
		WithExec([]string{
			"go", "install",
			"github.com/GeertJohan/fgt@latest",
		})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", install.File("/go/bin/golint")).
		WithFile("/", install.File("/go/bin/fgt"))
}
