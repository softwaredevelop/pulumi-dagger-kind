//revive:disable:package-comments
package main

import (
	"ci/util"
	"context"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

func main() {

	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c = c.Pipeline("cd")
	id, err := c.
		Container().
		// From("busybox:glibc").
		From("golang:alpine").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	if err != nil {
		panic(err)
	}

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..")
	if err != nil {
		panic(err)
	}

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	if err != nil {
		panic(err)
	}

	_, err = PulumiInstall(c, id).
		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
		WithExec([]string{"pulumi", "version"}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}
}
