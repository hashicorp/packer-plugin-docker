package docker

import (
	"context"
	"fmt"

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
	config, ok := state.Get("config").(*Config)
	if !ok {
		err := fmt.Errorf("error encountered obtaining docker config")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Building base image...")

	if config.EcrLogin {
		ui.Message("Fetching ECR credentials...")

		username, password, err := config.EcrGetLogin(config.LoginServer)
		if err != nil {
			err := fmt.Errorf("Error fetching ECR credentials: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		config.LoginUsername = username
		config.LoginPassword = password
	}

	if config.Login || config.EcrLogin {
		ui.Message("Logging in...")
		err := driver.Login(
			config.LoginServer,
			config.LoginUsername,
			config.LoginPassword)
		if err != nil {
			err := fmt.Errorf("Error logging in: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		defer func() {
			ui.Message("Logging out...")
			if err := driver.Logout(config.LoginServer); err != nil {
				ui.Error(fmt.Sprintf("Error logging out: %s", err))
			}
		}()
	}

	imageId, err := driver.Build(s.buildArgs.BuildArgs())
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Sayf("Finished building base image %q", imageId)

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
