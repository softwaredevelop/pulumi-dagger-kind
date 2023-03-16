//revive:disable:package-comments,exported
package util

import (
	"dagger.io/dagger"
)

func MountedHostDirectory(c *dagger.Client, id dagger.ContainerID, hostDir, mountedDir string) *dagger.Container {
	return mountedHostDirectory(c, id, hostDir, mountedDir)
}

func mountedHostDirectory(c *dagger.Client, id dagger.ContainerID, hostDir, mountedDir string) *dagger.Container {
	return c.Container(dagger.ContainerOpts{ID: id}).
		WithMountedDirectory(mountedDir, c.Host().Directory(hostDir, dagger.HostDirectoryOpts{
			Exclude: []string{".git"},
		}))
}
