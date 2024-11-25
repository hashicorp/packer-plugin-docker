// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config,AwsAccessConfig

package docker

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/mitchellh/mapstructure"
)

var (
	errArtifactNotUsed     = fmt.Errorf("No instructions given for handling the artifact; expected commit, discard, or export_path")
	errArtifactUseConflict = fmt.Errorf("Cannot specify more than one of commit, discard, and export_path")
	errExportPathNotFile   = fmt.Errorf("export_path must be a file, not a directory")
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`
	// Configuration for a bootstrap image derived from a Dockerfile
	//
	// Specifying this will make the builder run `docker build` on a provided
	// Dockerfile, and this image will then be used to perform the rest of
	// the build process.
	//
	// For more information on the contents of this object, refer to the
	// [Bootstrapping a build with a Dockerfile](#bootstrapping-a-build-with-a-dockerfile)
	// section of this documentation.
	BuildConfig DockerfileBootstrapConfig `mapstructure:"build"`
	// Set the author (e-mail) of a commit.
	Author string `mapstructure:"author"`
	// Dockerfile instructions to add to the commit. Example of instructions
	// are CMD, ENTRYPOINT, ENV, and EXPOSE. Example: [ "USER ubuntu", "WORKDIR
	// /app", "EXPOSE 8080" ]
	Changes []string `mapstructure:"changes"`
	// If true, the container will be committed to an image rather than exported.
	// Default `false`. If `commit` is `false`, then either `discard` must be
	// set to `true` or an `export_path` must be provided.
	Commit bool `mapstructure:"commit" required:"true"`
	// The directory inside container to mount temp directory from host server
	// for work [file provisioner](/packer/docs/provisioners/file). This defaults
	// to c:/packer-files on windows and /packer-files on other systems.
	ContainerDir string `mapstructure:"container_dir" required:"false"`
	// An array of devices which will be accessible in container when it's run
	// without `--privileged` flag.
	Device []string `mapstructure:"device" required:"false"`
	// Throw away the container when the build is complete. This is useful for
	// the [artifice
	// post-processor](/packer/docs/post-processors/artifice).
	Discard bool `mapstructure:"discard" required:"true"`
	// An array of additional [Linux
	// capabilities](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities)
	// to grant to the container.
	CapAdd []string `mapstructure:"cap_add" required:"false"`
	// An array of [Linux
	// capabilities](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities)
	// to drop from the container.
	CapDrop []string `mapstructure:"cap_drop" required:"false"`
	// Sets the docker binary to use for running commands.
	//
	// If you want to use a specific version of the docker binary, or a
	// docker alternative for building your container, you can specify this
	// through this option.
	// **Note**: if using an alternative like `podman`, not all options are
	// equivalent, and the build may fail in this case.
	//
	// Defaults to "docker"
	Executable string `mapstructure:"docker_path"`
	// Username (UID) to run remote commands with. You can also set the group
	// name/ID if you want: (UID or UID:GID). You may need this if you get
	// permission errors trying to run the shell or other provisioners.
	ExecUser string `mapstructure:"exec_user" required:"false"`
	// The path where the final container will be exported as a tar file.
	ExportPath string `mapstructure:"export_path" required:"true"`
	// The base image for the Docker container that will be started. This image
	// will be pulled from the Docker registry if it doesn't already exist.
	// Any value format that you can provide to `docker pull` is valid.
	// Example: `ubuntu` or `ubuntu:xenial`. If you only provide the repo, Docker
	// will pull the latest image, so setting `ubuntu` is the same as setting
	// `ubuntu:latest`. You can also set a distribution digest. For example,
	// ubuntu@sha256:a0d9e826ab87bd665cfc640598a871b748b4b70a01a4f3d174d4fb02adad07a9
	//
	// This cannot be used at the same time as `build`
	Image string `mapstructure:"image" required:"false"`
	// Set a message for the commit.
	Message string `mapstructure:"message" required:"true"`
	// If true, run the docker container with the `--privileged` flag. This
	// defaults to false if not set.
	Privileged bool `mapstructure:"privileged" required:"false"`
	Pty        bool
	// Set the container runtime. A runtime different from the one installed
	// by default with Docker (`runc`) must be installed and configured.
	// The possible values are (non-exhaustive list):
	// `runsc` for [gVisor](https://gvisor.dev/),
	// `kata-runtime` for [Kata Containers](https://katacontainers.io/),
	// `sysbox-runc` for [Nestybox](https://www.nestybox.com/).
	Runtime string `mapstructure:"runtime" required:"false"`
	// If true, the configured image will be pulled using `docker pull` prior
	// to use. Otherwise, it is assumed the image already exists and can be
	// used. This defaults to true if not set.
	//
	// If using `build`, this field will be ignored, as the `pull` option for
	// this operation will instead have precedence.
	Pull bool `mapstructure:"pull" required:"false"`
	// An array of arguments to pass to docker run in order to run the
	// container. By default this is set to `["-d", "-i", "-t",
	// "--entrypoint=/bin/sh", "--", "{{.Image}}"]` if you are using a linux
	// container, and `["-d", "-i", "-t", "--entrypoint=powershell", "--",
	// "{{.Image}}"]` if you are running a windows container. `{{.Image}}` is a
	// template variable that corresponds to the image template option. Passing
	// the entrypoint option this way will make it the default entrypoint of
	// the resulting image, so running docker run -it --rm  will start the
	// docker image from the /bin/sh shell interpreter; you could run a script
	// or another shell by running docker run -it --rm  -c /bin/bash. If your
	// docker image embeds a binary intended to be run often, you should
	// consider changing the default entrypoint to point to it.
	RunCommand []string `mapstructure:"run_command" required:"false"`
	// An array of additional tmpfs volumes to mount into this container.
	TmpFs []string `mapstructure:"tmpfs" required:"false"`
	// A mapping of additional volumes to mount into this container. The key of
	// the object is the host path, the value is the container path.
	Volumes map[string]string `mapstructure:"volumes" required:"false"`
	// If true, files uploaded to the container will be owned by the user the
	// container is running as. If false, the owner will depend on the version
	// of docker installed in the system. Defaults to true.
	FixUploadOwner bool `mapstructure:"fix_upload_owner" required:"false"`
	// If "true", tells Packer that you are building a Windows container
	// running on a windows host. This is necessary for building Windows
	// containers, because our normal docker bindings do not work for them.
	WindowsContainer bool `mapstructure:"windows_container" required:"false"`
	// Set platform if server is multi-platform capable
	//
	// This cannot be used at the same time as `build`; instead, use `build.platform`
	Platform string `mapstructure:"platform" required:"false"`

	// This is used to login to a private docker repository (e.g., dockerhub)
	// to build or pull a private base container. For pushing to a private
	//  repository, see the docker post-processors.
	Login bool `mapstructure:"login" required:"false"`
	// The password to use to authenticate to login.
	LoginPassword string `mapstructure:"login_password" required:"false"`
	// The server address to login to.
	LoginServer string `mapstructure:"login_server" required:"false"`
	// The username to use to authenticate to login.
	LoginUsername string `mapstructure:"login_username" required:"false"`
	// Defaults to false. If true, the builder will login in order to build or
	// pull the image from Amazon EC2 Container Registry (ECR). The builder
	// only logs in for the duration of the build or pull step. If true,
	// login_server is required and login, login_username, and login_password
	// will be ignored. For more information see the section on ECR.
	EcrLogin        bool `mapstructure:"ecr_login" required:"false"`
	AwsAccessConfig `mapstructure:",squash"`

	ctx interpolate.Context
}

func (c *Config) Prepare(raws ...interface{}) ([]string, error) {

	c.FixUploadOwner = true

	var md mapstructure.Metadata
	err := config.Decode(c, &config.DecodeOpts{
		Metadata:           &md,
		PluginType:         BuilderId,
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"run_command",
			},
		},
	}, raws...)
	if err != nil {
		return nil, err
	}

	// Defaults
	if len(c.RunCommand) == 0 {
		c.RunCommand = []string{"-d", "-i", "-t", "--entrypoint=/bin/sh", "--", "{{.Image}}"}
		if c.WindowsContainer {
			c.RunCommand = []string{"-d", "-i", "-t", "--entrypoint=powershell", "--", "{{.Image}}"}
		}
	}

	if c.Executable == "" {
		c.Executable = "docker"
	}

	// Default to the normal Docker type
	if c.Comm.Type == "" {
		c.Comm.Type = "docker"
		if c.WindowsContainer {
			c.Comm.Type = "dockerWindowsContainer"
		}
	}

	var errs *packersdk.MultiError
	var warnings []string

	if !c.BuildConfig.IsDefault() {
		_, err := c.BuildConfig.Prepare()
		if err != nil {
			errs = packersdk.MultiErrorAppend(errs, err)
		}

		if c.Image != "" {
			errs = packersdk.MultiErrorAppend(errs, errors.New("`image` cannot be specified with a build config"))
		}

		if c.Pull {
			warnings = append(warnings, "when running a bootstrap build, the `pull` option is ignored and is replaced by `build.pull` (true by default)")
			c.Pull = false
		}

		if c.Platform != "" {
			errs = packersdk.MultiErrorAppend(errs, errors.New("when running a bootstrap build, the `platform` option cannot be specified (use `build.platform` instead)"))
		}
		c.Platform = c.BuildConfig.Platform
	} else {
		// Default Pull if it wasn't set
		hasPull := false
		for _, k := range md.Keys {
			if k == "pull" {
				hasPull = true
				break
			}
		}

		if !hasPull {
			c.Pull = true
		}

		if c.Image == "" {
			errs = packersdk.MultiErrorAppend(errs,
				errors.New("missing 'image' attribute or 'build' section, either needs to be specified for a build to run."))
		}
	}

	if es := c.Comm.Prepare(&c.ctx); len(es) > 0 {
		errs = packersdk.MultiErrorAppend(errs, es...)
	}

	if (c.ExportPath != "" && c.Commit) || (c.ExportPath != "" && c.Discard) || (c.Commit && c.Discard) {
		errs = packersdk.MultiErrorAppend(errs, errArtifactUseConflict)
	}

	if c.ExportPath == "" && !c.Commit && !c.Discard {
		errs = packersdk.MultiErrorAppend(errs, errArtifactNotUsed)
	}

	if c.ExportPath != "" {
		if fi, err := os.Stat(c.ExportPath); err == nil && fi.IsDir() {
			errs = packersdk.MultiErrorAppend(errs, errExportPathNotFile)
		}
	}

	if c.ContainerDir == "" {
		if c.WindowsContainer {
			c.ContainerDir = "c:/packer-files"
		} else {
			c.ContainerDir = "/packer-files"
		}
	}

	if c.EcrLogin && c.LoginServer == "" {
		errs = packersdk.MultiErrorAppend(errs, fmt.Errorf("ECR login requires login server to be provided."))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return warnings, errs
	}

	return warnings, nil
}
