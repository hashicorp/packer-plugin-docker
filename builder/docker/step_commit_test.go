// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"
)

func testStepCommitState(t *testing.T) multistep.StateBag {
	state := testState(t)
	state.Put("container_id", "foo")
	return state
}

func TestStepCommit_impl(t *testing.T) {
	var _ multistep.Step = new(StepCommit)
}

func TestStepCommit(t *testing.T) {
	state := testStepCommitState(t)

	driver := state.Get("driver").(*MockDriver)
	driver.Sha256Result = "sha256:af61410def4ae2aece7c1b8d94b82ef434c8ee76e0e69001230f6636aea58cd1"
	driver.CommitImageId = "bar"

	step := &StepCommit{
		GeneratedData: &packerbuilderdata.GeneratedData{State: state},
	}
	defer step.Cleanup(state)

	// run the step
	if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %#v", action)
	}

	// verify we did the right thing
	if !driver.CommitCalled {
		t.Fatal("should've called")
	}

	// verify the ID is saved
	idRaw, ok := state.GetOk("image_id")
	if !ok {
		t.Fatal("should've saved ID")
	}

	id := idRaw.(string)
	if id != driver.CommitImageId {
		t.Fatalf("bad: %#v", id)
	}

	// Verify the sha256 is stored in generated data
	genData := state.Get("generated_data").(map[string]interface{})
	imSha := genData["ImageSha256"].(string)
	if imSha != driver.Sha256Result {
		t.Fatalf("Bad: image sha wasn't set properly; received %s", imSha)
	}
}

func TestStepCommit_error(t *testing.T) {
	state := testStepCommitState(t)
	step := new(StepCommit)
	defer step.Cleanup(state)

	driver := state.Get("driver").(*MockDriver)
	driver.CommitErr = errors.New("foo")

	// run the step
	if action := step.Run(context.Background(), state); action != multistep.ActionHalt {
		t.Fatalf("bad action: %#v", action)
	}

	// verify the ID is not saved
	if _, ok := state.GetOk("image_id"); ok {
		t.Fatal("shouldn't save image ID")
	}
}
