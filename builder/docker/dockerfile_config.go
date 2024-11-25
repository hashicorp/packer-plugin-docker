// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type DockerfileBootstrapConfig

package docker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/packer-plugin-sdk/template/config"
)

// DockerfileBootstrapConfig is the configuration for bootstrapping a docker
// builder with a user-provided Dockerfile.
//
// If used, the builder will build a source image locally using `docker build`,
// and continue the build normally with this as the base image.
type DockerfileBootstrapConfig struct {
	// Path to the dockerfile to use for building the base image
	//
	// If set, the builder will invoke `docker build` on it, and use the
	// produced image to continue the build afterwards.
	//
	// Note: Mutually exclusive with "image"
	DockerfilePath string `mapstructure:"path" required:"true"`
	// Directory to invoke `docker build` from
	//
	// Defaults to the directory from which we invoke packer.
	BuildDir string `mapstructure:"build_dir"`

	// A mapping of additional build args to provide. The key of
	// the object is the argument name, the value is the argument value.
	Arguments map[string]string `mapstructure:"arguments" required:"false"`

	// Set platform if server is multi-platform capable
	Platform string `mapstructure:"platform" required:"false"`

	// Pull the image when building the base docker image.
	//
	// Note: defaults to true, to disable this, explicitly set it to false.
	Pull config.Trilean `mapstructure:"pull"`
	// Compress the build context before sending to the docker daemon.
	//
	// This is especially useful if the build context is large, as copying it
	// can take a significant amount of time, while once compressed, this
	// can make builds faster, at the price of extra CPU resources.
	Compress bool `mapstructure:"compress"`
}

func (c *DockerfileBootstrapConfig) Prepare() ([]string, error) {
	// If the structure is all default, we can immediately return as we
	// won't invoke docker build to set a base image with it.
	if c.IsDefault() {
		return nil, nil
	}

	if c.BuildDir == "" {
		c.BuildDir = "."
	}

	st, err := os.Stat(c.DockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %q: %s", c.DockerfilePath, err)
	}

	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("dockerfile %q is not a regular file", c.DockerfilePath)
	}

	dockerfileAbsPath, err := filepath.Abs(c.DockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to compute absolute path for %q: %s", c.DockerfilePath, err)
	}

	c.DockerfilePath = dockerfileAbsPath

	st, err = os.Stat(c.BuildDir)
	if err != nil {
		return nil, fmt.Errorf("failed to stat build directory %q: %s", c.BuildDir, err)
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("specified build_dir %q is not a directory", c.BuildDir)
	}

	return nil, nil
}

// BuildArgs returns the list of arguments to pass to docker build.
func (c DockerfileBootstrapConfig) BuildArgs() []string {
	retArgs := []string{"-f", c.DockerfilePath}

	if c.Platform != "" {
		retArgs = append(retArgs, "--platform", c.Platform)
	}

	if !c.Pull.False() {
		retArgs = append(retArgs, "--pull")
	}

	if c.Compress {
		retArgs = append(retArgs, "--compress")
	}

	// Loops through map of build arguments to add to build command
	for key, value := range c.Arguments {
		arg := key + "=" + value
		retArgs = append(retArgs, "--build-arg", arg)
	}

	return append(retArgs, c.BuildDir)
}

func (c DockerfileBootstrapConfig) IsDefault() bool {
	if c.BuildDir != "" {
		return false
	}

	if c.Compress {
		return false
	}

	if c.DockerfilePath != "" {
		return false
	}

	return true
}
