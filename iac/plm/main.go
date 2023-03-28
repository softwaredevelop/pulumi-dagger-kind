//revive:disable:package-comments
package main

import (
	"context"
	"os"
	"path/filepath"
	"plm/util"

	"dagger.io/dagger"
)

func main() {

	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c = c.Pipeline("cd-inline-source")
	// c = c.Pipeline("cd-local-source")
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
	p := filepath.Join(dir, "..", "ils")
	// p := filepath.Join(dir, "..", "lcs")
	if err != nil {
		panic(err)
	}

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	if err != nil {
		panic(err)
	}

	pat := os.Getenv("PULUMI_ACCESS_TOKEN")
	ght := os.Getenv("GITHUB_TOKEN")
	gho := os.Getenv("GITHUB_OWNER")
	id, err = util.PulumiInstall(c, id).
		Pipeline("pulumi").
		WithWorkdir(mountedDir).
		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
		WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
		WithEnvVariable("PULUMI_ACCESS_TOKEN", pat).
		WithEnvVariable("GITHUB_TOKEN", ght).
		WithEnvVariable("GITHUB_OWNER", gho).
		WithExec([]string{"pulumi", "login"}).
		ID(ctx)
	if err != nil {
		panic(err)
	}

	// args := os.Args[1:]
	// if len(args) > 0 {
	// 	if args[0] == "destroy" {
	// 		_, err = c.Container(dagger.ContainerOpts{ID: id}).
	// 			Pipeline("pulumi-inline-source3").
	// 			WithWorkdir(mountedDir).
	// 			WithExec([]string{"go", "run", "-v", "inline.go", "destroy"}).
	// 			Stdout(ctx)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }

	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-inline-source2").
		WithWorkdir(mountedDir).
		WithExec([]string{"go", "run", "-v", "inline.go"}).
		// WithExec([]string{"go", "run", "-v", "local.go"}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}
}
