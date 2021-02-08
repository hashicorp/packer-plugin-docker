package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-docker/builder/docker"
	dockerimport "github.com/hashicorp/packer-plugin-docker/post-processor/docker-import"
	dockerpush "github.com/hashicorp/packer-plugin-docker/post-processor/docker-push"
	dockersave "github.com/hashicorp/packer-plugin-docker/post-processor/docker-save"
	dockertag "github.com/hashicorp/packer-plugin-docker/post-processor/docker-tag"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"
)

var (
	// Version is the main version number that is being run at the moment.
	Version = "0.0.1"

	// VersionPrerelease is A pre-release marker for the Version. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this
	// is a pre-release such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "dev"

	// PluginVersion is used by the plugin set to allow Packer to recognize
	// what version this plugin is.
	PluginVersion = version.InitializePluginVersion(Version, VersionPrerelease)
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(docker.Builder))
	pps.RegisterPostProcessor("import", new(dockerimport.PostProcessor))
	pps.RegisterPostProcessor("push", new(dockerpush.PostProcessor))
	pps.RegisterPostProcessor("save", new(dockersave.PostProcessor))
	pps.RegisterPostProcessor("tag", new(dockertag.PostProcessor))
	pps.SetVersion(PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
