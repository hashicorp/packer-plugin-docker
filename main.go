// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-docker/builder/docker"
	dockerimport "github.com/hashicorp/packer-plugin-docker/post-processor/docker-import"
	dockerpush "github.com/hashicorp/packer-plugin-docker/post-processor/docker-push"
	dockersave "github.com/hashicorp/packer-plugin-docker/post-processor/docker-save"
	dockertag "github.com/hashicorp/packer-plugin-docker/post-processor/docker-tag"
	"github.com/hashicorp/packer-plugin-docker/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(docker.Builder))
	pps.RegisterPostProcessor("import", new(dockerimport.PostProcessor))
	pps.RegisterPostProcessor("push", new(dockerpush.PostProcessor))
	pps.RegisterPostProcessor("save", new(dockersave.PostProcessor))
	pps.RegisterPostProcessor("tag", new(dockertag.PostProcessor))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
