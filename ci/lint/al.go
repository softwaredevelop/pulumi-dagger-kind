//revive:disable:package-comments,exported
package lint

import "dagger.io/dagger"

func Al(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return al(c, id)
}

func al(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("golang:alpine").
		WithExec([]string{"go", "install",
			"github.com/rhysd/actionlint/cmd/actionlint@latest",
		})

	shellcheck := c.
		Container().
		From("koalaman/shellcheck-alpine:stable")

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/usr/bin", install.File("/go/bin/actionlint")).
		WithFile("/usr/bin", shellcheck.File("/bin/shellcheck"))
}
