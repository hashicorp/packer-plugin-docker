// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepBuild struct {
	buildArgs DockerfileBootstrapConfig
	ran       bool
}

func (s *stepBuild) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if s.buildArgs.IsDefault() {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packer.Ui)
	driver := state.Get("driver").(Driver)
	ui.Say("Building base image...")

	imageId, err := driver.Build(s.buildArgs.BuildArgs())
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Sayf("Finished building base image %q", imageId)

	cfg, ok := state.GetOk("config")
	if !ok {
		state.Put("error", errors.New("missing config in state; this is a docker plugin bug, please report upstream"))
		return multistep.ActionHalt
	}

	config, ok := cfg.(*Config)
	if !ok {
		state.Put("error", fmt.Errorf("config object set but type (%s) doesn't match. This is a docker plugin bug, please report upstream", reflect.TypeOf(cfg).String()))
		return multistep.ActionHalt
	}

	config.Image = imageId

	s.ran = true

	return multistep.ActionContinue
}

func (s *stepBuild) Cleanup(state multistep.StateBag) {
	if !s.ran {
		return
	}

	ui := state.Get("ui").(packer.Ui)
	config, ok := state.Get("config").(*Config)
	if !ok {
		ui.Say("missing config object in state, won't cleanup built image")
		return
	}

	if !config.Discard {
		ui.Say("final image is not discarded, removing the built image will fail because of existing dependencies, skipping cleanup for docker build.")
		return
	}

	driver := state.Get("driver").(Driver)

	err := driver.DeleteImage(config.Image)
	if err != nil {
		ui.Sayf("failed to remove image %q: %s", config.Image, err)
		ui.Say("if you have other images using this dockerfile, this is expected and can safely be ignored.")
	}
}
