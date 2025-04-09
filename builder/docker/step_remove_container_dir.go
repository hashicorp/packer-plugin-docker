package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepRemoveContainerDir struct {
	ContainerDir string
}

func (s *StepRemoveContainerDir) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {

	comm := state.Get("communicator").(packersdk.Communicator)
	ui := state.Get("ui").(packersdk.Ui)

	cmd := new(packersdk.RemoteCmd)

	ui.Sayf("Trying to remove temporary container directory that was mounted. directory: %q", s.ContainerDir)

	cmd.Command = fmt.Sprintf("whoami; umount %s; rm -rf %s;", s.ContainerDir, s.ContainerDir)
	if err := cmd.RunWithUi(ctx, comm, ui); err != nil {
		log.Printf("Error cleaning temporary directory: %s", err)
	}

	return multistep.ActionContinue

}

func (s *StepRemoveContainerDir) Cleanup(state multistep.StateBag) {}
