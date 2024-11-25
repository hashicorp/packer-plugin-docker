// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

var simpleDockerTemplate = `
source "docker" "simple" {
	image = "alpine:latest"
	discard = true
}

build {
	sources = ["docker.simple"]
}
`

func TestAccBuilder_Basic(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "docker-builder-simple-case",
		Template: simpleDockerTemplate,
		Check: func(c *exec.Cmd, s string) error {
			if c.ProcessState.ExitCode() != 0 {
				return fmt.Errorf("build failed unexpectedly with exit code %d", c.ProcessState.ExitCode())
			}
			return nil
		},
	}

	acctest.TestPlugin(t, testCase)
}

var dockerBuildTemplate = `
source "docker" "with_build" {
	build {
		path = "./test-fixtures/sample_dockerfile"
	}
	discard = true
}

build {
	sources = ["docker.with_build"]

	provisioner "shell" {
		inline = ["echo \"alpine release is $(cat /etc/alpine-release)\""]
	}
}
`

func TestAccBuilder_DockerBuildSimple(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "docker-builder-with-build-dockerfile",
		Template: dockerBuildTemplate,
		Check: func(c *exec.Cmd, logfile string) error {
			if c.ProcessState.ExitCode() != 0 {
				return fmt.Errorf("build failed unexpectedly with exit code %d", c.ProcessState.ExitCode())
			}

			logs, err := os.ReadFile(logfile)
			if err != nil {
				return fmt.Errorf("failed to read logs of packer build %q: %s", logfile, err)
			}

			re := regexp.MustCompile(`alpine release is ([0-9]+\.){2}[0-9]+`)
			if !re.Match(logs) {
				return fmt.Errorf("alpine release is not present on the image, should have been.")
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
