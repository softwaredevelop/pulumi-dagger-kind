//revive:disable:package-comments
package main

import (
	"ci/lint"
	"ci/util"
	"context"
	"log"
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

	c = c.Pipeline("ci")
	id, err := c.
		Container().
		From("busybox:uclibc").
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

	ec, err := lint.Ec(c, id).
		Pipeline("ec").
		WithWorkdir(mountedDir).
		WithExec([]string{"editorconfig-checker", "-verbose"}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(ec)

	al, err := lint.Al(c, id).
		Pipeline("actionlint").
		WithWorkdir(mountedDir).
		WithExec([]string{"actionlint", "-pyflakes=", "-verbose", ".github/workflows/dagger.yml"}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(al)

	revive, err := lint.Revive(c, id).
		Pipeline("revive").
		WithWorkdir(mountedDir).
		WithExec([]string{"revive", "-set_exit_status", "./..."}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(revive)
}
