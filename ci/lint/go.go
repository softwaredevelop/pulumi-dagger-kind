//revive:disable:package-comments,exported
package lint

import "dagger.io/dagger"

func Revive(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return revive(c, id)
}

func revive(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("golang:alpine").
		WithExec([]string{"go", "install",
			"github.com/mgechev/revive@latest",
		})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/usr/bin", install.File("/go/bin/revive"))
}

func GoLint(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return goLint(c, id)
}

func goLint(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("golang:alpine").
		WithExec([]string{"go", "install",
			"golang.org/x/lint/golint@latest",
		}).
		WithExec([]string{"go", "install",
			"github.com/GeertJohan/fgt@latest",
		})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/usr/bin", install.File("/go/bin/golint")).
		WithFile("/usr/bin", install.File("/go/bin/fgt"))
}
