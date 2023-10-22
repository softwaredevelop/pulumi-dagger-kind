//revive:disable:package-comments
package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"cc/format"
	"cc/linting"
	"cc/util"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c = c.Pipeline("code_quality")
	id, err := c.
		Container().
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

	var wg sync.WaitGroup
	wg.Add(4)

	runtime.GOMAXPROCS(runtime.NumCPU())

	go func() {
		defer wg.Done()
		f := c.Pipeline("gofumpt")
		err = format.Gofumpt(f, id, mountedDir)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		e := c.Pipeline("editorconfig_checker")
		err = linting.EditorconfigChecker(e, id, mountedDir)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		a := c.Pipeline("actionlint")
		err = linting.Actionlint(dir, a, id, mountedDir)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		r := c.Pipeline("revive")
		err = linting.Revive(r, id, mountedDir)
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
