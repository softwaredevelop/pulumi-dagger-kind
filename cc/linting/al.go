//revive:disable:package-comments,exported
package linting

import (
	"cc/util"
	"context"
	"path/filepath"

	"dagger.io/dagger"
)

func Actionlint(dir string, c *dagger.Client, id dagger.ContainerID, mountedDir string) error {
	ctx := context.Background()
	p := filepath.Join(dir, "..", ".github", "workflows")
	id, err := util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	if err != nil {
		return err
	}
	_, err = al(c, id).
		WithWorkdir(mountedDir).
		WithExec([]string{"/actionlint",
			"-debug",
			"-pyflakes",
			"-shellcheck",
			"-verbose",
			"code-quality.yml",
		}).
		WithExec([]string{"/actionlint",
			"-debug",
			"-pyflakes",
			"-shellcheck",
			"-verbose",
			"iac.yml",
		}).
		WithExec([]string{"/actionlint",
			"-debug",
			"-pyflakes",
			"-shellcheck",
			"-verbose",
			"test.yml",
		}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Al(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return al(c, id)
}

func al(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("golang:alpine").
		WithExec([]string{
			"go", "install",
			"github.com/rhysd/actionlint/cmd/actionlint@latest",
		})

	shellcheck := c.
		Container().
		From("koalaman/shellcheck-alpine:stable")

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", install.File("/go/bin/actionlint")).
		WithFile("/", shellcheck.File("/bin/shellcheck"))
}
