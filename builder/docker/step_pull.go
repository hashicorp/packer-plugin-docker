package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"
)

type StepPull struct {
	GeneratedData *packerbuilderdata.GeneratedData
}

func (s *StepPull) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)
	config, ok := state.Get("config").(*Config)
	if !ok {
		err := fmt.Errorf("error encountered obtaining docker config")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if !config.Pull {
		log.Println("Pull disabled, won't docker pull")
		return multistep.ActionContinue
	}

	ui.Say(fmt.Sprintf("Pulling Docker image: %s", config.Image))

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

	driver := state.Get("driver").(Driver)
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

	if err := driver.Pull(config.Image); err != nil {
		err := fmt.Errorf("Error pulling Docker image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Store data about source image.
	// Distribution digest is something you can use to pull the image down.
	sourceDigest, err := driver.Digest(config.Image)
	if err != nil {
		err := fmt.Errorf("Error determining source Docker image digest: %s", err)
		ui.Error(err.Error())
	}
	state.Put("source_digest", sourceDigest)
	s.GeneratedData.Put("SourceImageDigest", sourceDigest)

	// Image Id is a shasum that is unique to this image.
	sourceSha256, err := driver.Sha256(config.Image)
	if err != nil {
		err := fmt.Errorf("Error determining source Docker image Id: %s", err)
		ui.Error(err.Error())
	}
	state.Put("source_sha256", sourceSha256)
	s.GeneratedData.Put("SourceImageSha256", sourceSha256)

	return multistep.ActionContinue
}

func (s *StepPull) Cleanup(state multistep.StateBag) {
}
