package main

import (
	"context"
	"os"
	"path/filepath"
	"plm/util"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestPulumiInlineSourceService(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	c = c.Pipeline("cd-inline-source-service-test")
	require.NotNil(t, c)

	id, err := c.
		Container().
		// From("busybox:glibc").
		From("golang:alpine").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..", "ils")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	pat := os.Getenv("PULUMI_ACCESS_TOKEN")
	pon := os.Getenv("PULUMI_ORG_NAME")
	id, err = util.PulumiInstall(c, id).
		Pipeline("pulumi").
		WithWorkdir(mountedDir).
		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
		WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
		WithEnvVariable("PULUMI_ACCESS_TOKEN", pat).
		WithEnvVariable("PULUMI_ORG_NAME", pon).
		WithExec([]string{"pulumi", "login"}).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	reMatching := "TestUpsertStackInlineSourceRefresh$"
	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-inline-source-upster-test1").
		WithWorkdir(mountedDir).
		WithExec([]string{"go", "test", "-v", "-run", reMatching}).
		Stdout(ctx)
	require.NoError(t, err)
}

func TestPulumiInlineSource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	c = c.Pipeline("cd-inline-source-test")
	require.NotNil(t, c)

	id, err := c.
		Container().
		// From("busybox:glibc").
		From("golang:alpine").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..", "ils")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	ght := os.Getenv("GITHUB_TOKEN")
	gho := os.Getenv("GITHUB_OWNER")
	id, err = util.PulumiInstall(c, id).
		Pipeline("pulumi").
		WithWorkdir(mountedDir).
		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
		WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
		WithEnvVariable("GITHUB_TOKEN", ght).
		WithEnvVariable("GITHUB_OWNER", gho).
		WithExec([]string{"pulumi", "login", "--local"}).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// reMatching := "TestUpsertStackInlineSource$"
	// _, err = c.Container(dagger.ContainerOpts{ID: id}).
	// 	Pipeline("pulumi-inline-source-test1").
	// 	WithWorkdir(mountedDir).
	// 	WithExec([]string{"go", "test", "-v", "inline_test.go", "-run", reMatching}).
	// 	Stdout(ctx)
	// require.NoError(t, err)

	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-inline-source-test2").
		WithWorkdir(mountedDir).
		WithExec([]string{"go", "test", "-v", "inline_test.go"}).
		Stdout(ctx)
	require.NoError(t, err)
}

func TestPulumiLocalSource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	c = c.Pipeline("cd-local-source-test")
	require.NotNil(t, c)

	id, err := c.
		Container().
		// From("busybox:glibc").
		From("golang:alpine").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..", "lcs")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	id, err = util.PulumiInstall(c, id).
		Pipeline("pulumi").
		WithWorkdir(mountedDir).
		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
		WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
		WithExec([]string{"pulumi", "login", "--local"}).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	reMatching := "TestNewStackLocalSourceWorkspaceEnvVars$"
	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-local-source1-test").
		WithWorkdir(mountedDir).
		WithExec([]string{"go", "test", "-v", "-run", reMatching}).
		Stdout(ctx)
	require.NoError(t, err)

	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi-local-source2-test").
		WithWorkdir(mountedDir).
		WithExec([]string{"go", "test", "-v", "local_test.go"}).
		Stdout(ctx)
	require.NoError(t, err)
}

// func TestPulumiLocalSourceHost(t *testing.T) {
// 	t.Parallel()
// 	ctx := context.Background()

// 	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
// 	require.NoError(t, err)
// 	defer c.Close()

// 	c = c.Pipeline("cd")
// 	require.NotNil(t, c)

// 	id, err := c.
// 		Container().
// 		// From("busybox:glibc").
// 		From("golang:alpine").
// 		WithMountedTemp("/mountedtmp").
// 		ID(ctx)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, id)

// 	dir, _ := os.Getwd()
// 	p := filepath.Join(dir, "..", "lsh")
// 	require.NoError(t, err)
// 	require.NotEmpty(t, p)

// 	mountedDir := "/mountedtmp"
// 	id, err = MountedHostDirectory(c, id, p, mountedDir).
// 		ID(ctx)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, id)

// 	id, err = PulumiInstall(c, id).
// 		Pipeline("pulumi").
// 		WithWorkdir(mountedDir).
// 		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
// 		WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
// 		WithExec([]string{"pulumi", "login", "--local"}).
// 		ID(ctx)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, id)

// 	_, err = c.Container(dagger.ContainerOpts{ID: id}).
// 		Pipeline("pulumi-local-source-host").
// 		WithWorkdir(mountedDir).
// 		WithExec([]string{"pulumi", "stack",
// 			"init",
// 			"--stack=lsh",
// 			"--color=always",
// 			"--non-interactive",
// 		}).
// 		WithExec([]string{"pulumi", "stack",
// 			"ls",
// 			"--color=always",
// 			"--non-interactive",
// 		}).
// 		WithExec([]string{"pulumi", "preview",
// 			"--debug",
// 			"--color=always",
// 			"--non-interactive",
// 		}).
// 		Stdout(ctx)
// 	require.NoError(t, err)

// 	_, err = c.Container(dagger.ContainerOpts{ID: id}).
// 		Pipeline("pulumi-local-source-host").
// 		WithWorkdir(mountedDir).
// 		WithExec([]string{"pulumi", "stack",
// 			"init",
// 			"--stack=lsh",
// 			"--color=always",
// 			"--non-interactive",
// 		}).
// 		WithExec([]string{"pulumi", "up",
// 			"--debug",
// 			"--skip-preview",
// 			"--color=always",
// 			"--non-interactive",
// 		}).
// 		Stdout(ctx)
// 	require.NoError(t, err)
// }

func TestMountedHostParentDirectory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	container := c.Container().From("busybox:glibc")
	require.NotNil(t, container)

	dir, _ := os.Getwd()
	p := filepath.Join(dir, "..")
	require.NoError(t, err)
	require.NotEmpty(t, p)
	id, err := container.
		WithMountedDirectory("/mountedtmp", c.Host().Directory(p)).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	mntdir, err := c.
		Container(dagger.ContainerOpts{ID: id}).
		WithExec([]string{"ls", "/mountedtmp"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, mntdir)

	lsdir, err := c.
		Container(dagger.ContainerOpts{ID: id}).
		WithWorkdir("/mountedtmp").
		WithExec([]string{"ls", "-la"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, lsdir)
}

func TestContainerID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	container := c.Container().From("alpine")
	require.NotNil(t, container)

	id, err := container.ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	releaseName, err := c.Container(dagger.ContainerOpts{ID: id}).
		WithExec([]string{"/bin/sh", "-c", "cat /etc/os-release | awk -F= '/^NAME/ {print $2}' | tr -d '\"'"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "Alpine Linux\n", releaseName)
}

func TestErrorMessage(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	_, err = c.
		Container().
		From("fake.invalid").
		ID(ctx)
	require.Error(t, err)
	require.ErrorContains(t, err, "not exist")
}

func TestConnect(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()
}

func TestPulumiInstallLogin(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	container := c.Container().From("busybox:uclibc")
	require.NotNil(t, container)

	id, err := container.ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	container = c.Container(dagger.ContainerOpts{ID: id})
	require.NotNil(t, container)

	container = util.PulumiInstall(c, id)
	require.NotNil(t, container)
	out, err := container.
		WithExec([]string{"ls", "/usr/bin/pulumi"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/pulumi\n", out)

	out, err = container.
		WithExec([]string{"pulumi", "login", "--local"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func TestPulumiGoInstall(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	container := c.Container().From("busybox:glibc")
	require.NotNil(t, container)

	id, err := container.ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	container = util.PulumiGoInstall(c, id)
	require.NotNil(t, container)

	out, err := container.
		WithExec([]string{"go", "version"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Contains(t, out, "version")

	out, err = container.
		WithExec([]string{"pulumi", "version"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Contains(t, out, "v")
}

func TestPulumiOnPath(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	container := c.Container().From("busybox:uclibc")
	require.NotNil(t, container)

	id, err := container.ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	container = util.PulumiInstall(c, id)
	require.NotNil(t, container)

	_, err = container.
		WithExec([]string{"pulumi", "login", "--local"}).
		Stdout(ctx)
	require.NoError(t, err)
}

func TestEnvVariable(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	container := c.Container().From("busybox:glibc")
	require.NotNil(t, container)

	id, err := container.ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	container = c.Container(dagger.ContainerOpts{ID: id})
	require.NotNil(t, container)

	env, err := container.EnvVariable(ctx, "PATH")
	require.NoError(t, err)
	require.Contains(t, env, "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")

	env, err = container.
		WithEnvVariable("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin").
		EnvVariable(ctx, "PATH")
	require.NoError(t, err)
	require.Contains(t, env, "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
}

func TestPulumiInstall(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	container := c.Container().From("busybox:uclibc")
	require.NotNil(t, container)

	id, err := container.ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	container = c.Container(dagger.ContainerOpts{ID: id})
	require.NotNil(t, container)

	container = util.PulumiInstall(c, id)
	require.NotNil(t, container)
	out, err := container.
		WithExec([]string{"ls", "/usr/bin/pulumi"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/pulumi\n", out)
}
