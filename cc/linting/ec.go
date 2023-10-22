//revive:disable:package-comments,exported
package linting

import (
	"context"
	"log"

	"dagger.io/dagger"
)

func EditorconfigChecker(c *dagger.Client, id dagger.ContainerID, mountedDir string) error {
	ctx := context.Background()
	_, err := ec(c, id).
		WithWorkdir(mountedDir).
		WithExec([]string{"/editorconfig-checker", "-verbose"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Ec(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return ec(c, id)
}

func ec(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("golang:alpine").
		WithExec([]string{
			"go", "install",
			"github.com/editorconfig-checker/editorconfig-checker/cmd/editorconfig-checker@latest",
		})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", install.File("/go/bin/editorconfig-checker"))
}

func EcVersion() (string, error) {
	return ecVersion()
}

func ecVersion() (string, error) {
	ctx := context.Background()
	c, err := dagger.Connect(ctx)
	if err != nil {
		log.Println(err)
		return "", err
	}
	repoURL := "https://github.com/editorconfig-checker/editorconfig-checker.git"
	fileName := "VERSION"
	version, err := gitCloneFileContent(ctx, c, repoURL, fileName)
	if err != nil {
		return "", err
	}
	return version, nil
}

func Ec2(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	repoURL := "https://github.com/editorconfig-checker/editorconfig-checker.git"
	version, err := ecVersion()
	if err != nil {
		log.Println(err)
		return nil
	}
	build := c.
		Container().
		From("alpine").
		WithExec([]string{
			"apk", "add",
			"--no-cache",
			"go",
			"git",
		}).
		WithWorkdir("/ec").
		WithExec([]string{
			"git", "clone",
			"--single-branch",
			"--branch", "main",
			repoURL,
			"/ec",
		}).
		WithEnvVariable("GO111MODULE", "on").
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{
			"go", "build",
			"-ldflags", "-X main.version=" + version,
			"-o", "bin/ec",
			"./cmd/editorconfig-checker/main.go",
		})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", build.File("/ec/bin/ec"))
}

func Ec1(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return ec1(c, id)
}

func ec1(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	repoURL := "https://github.com/editorconfig-checker/editorconfig-checker.git"
	version, err := ecVersion()
	if err != nil {
		log.Println(err)
		return nil
	}
	build := c.
		Container().
		From("golang:alpine").
		WithExec([]string{
			"apk", "add",
			"--no-cache",
			"git",
		}).
		WithWorkdir("/ec").
		WithExec([]string{
			"git", "clone",
			"--single-branch",
			"--branch", "main",
			repoURL,
			"/ec",
		}).
		WithEnvVariable("GO111MODULE", "on").
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{
			"go", "build",
			"-ldflags", "-X main.version=" + version,
			"-o", "bin/ec",
			"./cmd/editorconfig-checker/main.go",
		})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", build.File("/ec/bin/ec"))
}

func GitCloneFileContent(ctx context.Context, c *dagger.Client, repoURL, fileName string) (string, error) {
	return gitCloneFileContent(ctx, c, repoURL, fileName)
}

func gitCloneFileContent(ctx context.Context, c *dagger.Client, repoURL, fileName string) (string, error) {
	contents, err := gitClone(c, repoURL).
		File(fileName).Contents(ctx)
	if err != nil {
		return "", err
	}
	return contents, nil
}

func GitClone(c *dagger.Client, repoURL string) *dagger.Directory {
	return gitClone(c, repoURL)
}

func gitClone(c *dagger.Client, repoURL string) *dagger.Directory {
	return c.
		Git(repoURL, dagger.GitOpts{KeepGitDir: true}).
		Branch("main").
		Tree()
}
