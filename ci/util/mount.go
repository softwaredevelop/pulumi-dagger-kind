//revive:disable:package-comments,exported
package util

import (
	"os"
	"path"

	"dagger.io/dagger"
)

func HostDirectoryPath(e string) (string, error) {
	return hostDirectoryPath(e)
}

func hostDirectoryPath(e string) (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(p, e), nil
}

func MountedHostDirectory(c *dagger.Client, id dagger.ContainerID, hostDir, mountedDir string) *dagger.Container {
	return mountedHostDirectory(c, id, hostDir, mountedDir)
}

func mountedHostDirectory(c *dagger.Client, id dagger.ContainerID, hostDir, mountedDir string) *dagger.Container {
	return c.Container(dagger.ContainerOpts{ID: id}).
		WithMountedDirectory(mountedDir, c.Host().Directory(hostDir, dagger.HostDirectoryOpts{
			Exclude: []string{".git"},
		}))
}
