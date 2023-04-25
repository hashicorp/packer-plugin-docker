// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"
)

// StepCommit commits the container to a image.
type StepCommit struct {
	imageId       string
	GeneratedData *packerbuilderdata.GeneratedData
}

func (s *StepCommit) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)
	config, ok := state.Get("config").(*Config)
	if !ok {
		err := fmt.Errorf("error encountered obtaining docker config")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	driver := state.Get("driver").(Driver)
	containerId := state.Get("container_id").(string)
	if config.WindowsContainer {
		// docker can't commit a running Windows container
		err := driver.StopContainer(containerId)
		if err != nil {
			state.Put("error", err)
			ui.Error(fmt.Sprintf("Error halting windows container for commit: %s",
				err.Error()))
			return multistep.ActionHalt
		}
	}
	ui.Say("Committing the container")
	imageId, err := driver.Commit(containerId, config.Author, config.Changes, config.Message)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Save the container ID to state and to generated data
	s.imageId = imageId
	state.Put("image_id", s.imageId)
	s256, err := driver.Sha256(s.imageId)
	if err == nil {
		s.GeneratedData.Put("ImageSha256", s256)
	}

	ui.Message(fmt.Sprintf("Image ID: %s", s.imageId))

	return multistep.ActionContinue
}

func (s *StepCommit) Cleanup(state multistep.StateBag) {}
