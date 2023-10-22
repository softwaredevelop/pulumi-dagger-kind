//revive:disable:package-comments,exported
package linting

import (
	"cc/util"
	"context"

	"dagger.io/dagger"
)

func Yamllint(p string, c *dagger.Client) error {
	ctx := context.Background()

	id, err := c.
		Container().
		From("pipelinecomponents/yamllint").
		WithMountedTemp("/code").
		ID(ctx)
	if err != nil {
		return err
	}

	mountedDir := "/code"
	_, err = util.MountedHostDirectory(c, id, p, mountedDir).
		WithWorkdir(mountedDir).
		WithExec([]string{"yamllint",
			"--config-data",
			"{extends: default, rules: {line-length: {level: warning}}}",
			"--no-warnings",
			"."}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}
