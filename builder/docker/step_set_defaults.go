// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"context"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type StepSetDefaults struct{}

func (s *StepSetDefaults) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(*DockerDriver)
	config := state.Get("config").(*Config)

	// Fetch default CMD and ENTRYPOINT
	//
	// NOTE: if they're empty, the returned value is expected to be [""]
	// This is because in order to override what we had set when running
	// the build container, they need to be an array with nothing in it.
	// The commit command treats explicit `null` as a NOOP, same with `[]`.
	// If a string is passed as the change argument, it will be appended to
	// what already exists for the cmd/entrypoint, which doesn't match the
	// need here, as we need to forcefully set what was originally
	// specified.
	//
	// So while not necessarily clean, the best way to force a similar
	// behaviour as the original image, we default on an array with an
	// empty string as argument, which is effectively the same as `null`.
	defaultCmd, _ := driver.Cmd(config.Image)
	defaultEntrypoint, _ := driver.Entrypoint(config.Image)

	// Set defaults if not provided by the user
	hasCmd, hasEntrypoint := false, false
	for _, change := range config.Changes {
		if strings.HasPrefix(change, "CMD") {
			hasCmd = true
		} else if strings.HasPrefix(change, "ENTRYPOINT") {
			hasEntrypoint = true
		}
	}

	if !hasCmd {
		config.Changes = append(config.Changes, "CMD "+defaultCmd)
	}
	if !hasEntrypoint {
		config.Changes = append(config.Changes, "ENTRYPOINT "+defaultEntrypoint)
	}

	return multistep.ActionContinue
}

func (s *StepSetDefaults) Cleanup(state multistep.StateBag) {}
