//revive:disable:package-comments,exported
package util

import (
	"runtime"
	"strings"

	"dagger.io/dagger"
)

func P(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return pulumiScratchInstall(c, id)
}

func pulumiScratchInstall(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	pulumiInstall := c.
		Container().
		From("alpine").
		WithEnvVariable("HOME", "/tmp").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add",
			"--no-cache",
			"curl",
		}).
		WithExec([]string{"/bin/sh", "-c", "curl -fsSL https://get.pulumi.com | sh"})
	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/pulumi", pulumiInstall.File("/tmp/.pulumi/bin/pulumi")).
		WithFile("/pulumi-language-go", pulumiInstall.File("/tmp/.pulumi/bin/pulumi-language-go"))
}

func PulumiGoInstall(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return pulumiGoInstall(c, id)
}

func pulumiGoInstall(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	pulumiInstall := c.
		Container().
		From("alpine").WithEnvVariable("HOME", "/tmp").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add",
			"--no-cache",
			"curl",
		}).
		WithExec([]string{"/bin/sh", "-c", "curl -fsSL https://get.pulumi.com | sh"})

	version := strings.TrimPrefix(runtime.Version(), "go")
	goDownload := c.
		Container().
		From("busybox:glibc").
		// From("alpine").
		WithExec([]string{"wget", "https://go.dev/dl/go" + version + ".linux-amd64.tar.gz"}).
		WithExec([]string{"mkdir", "-p", "/usr/local/go/bin"}).
		WithExec([]string{"tar", "-xzf",
			"go" + version + ".linux-amd64.tar.gz",
			"--strip-components=2",
			"go/bin/go",
			"-C",
			"/usr/local/go/bin/",
		})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/usr/bin", pulumiInstall.File("/tmp/.pulumi/bin/pulumi")).
		WithFile("/usr/bin", pulumiInstall.File("/tmp/.pulumi/bin/pulumi-language-go")).
		WithExec([]string{"mkdir", "-p", "/usr/local/go/bin"}).
		WithFile("/usr/local/go/bin", goDownload.File("/usr/local/go/bin/go")).
		WithEnvVariable("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin")
}

func PulumiInstall(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return pulumiInstall(c, id)
}

func pulumiInstall(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("alpine").WithEnvVariable("HOME", "/tmp").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add",
			"--no-cache",
			"curl",
		}).
		WithExec([]string{"/bin/sh", "-c", "curl -fsSL https://get.pulumi.com | sh"})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/usr/bin", install.File("/tmp/.pulumi/bin/pulumi")).
		WithFile("/usr/bin", install.File("/tmp/.pulumi/bin/pulumi-language-go"))
}
