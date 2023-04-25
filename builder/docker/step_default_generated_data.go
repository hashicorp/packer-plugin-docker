// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"
)

// StepDefaultGeneratedData Adds the placeholders for special generated values
// that Docker is expected to return. This makes sure that the accessor has
// _something_ to read in the provisioners, regardless of whether the value.
// was created.
// The true values are put in generated data in the steps where the
// values are actually created (step pull, and step commit)
type StepDefaultGeneratedData struct {
	GeneratedData *packerbuilderdata.GeneratedData
}

func (s *StepDefaultGeneratedData) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	s.GeneratedData.Put("ImageSha256", "ERR_IMAGE_SHA256_NOT_FOUND")
	s.GeneratedData.Put("SourceImageDigest", "ERR_SOURCE_IMAGE_DIGEST_NOT_FOUND")
	s.GeneratedData.Put("SourceImageSha256", "ERR_SOURCE_IMAGE_SHA256_NOT_FOUND")

	return multistep.ActionContinue
}

func (s *StepDefaultGeneratedData) Cleanup(_ multistep.StateBag) {
	// No cleanup...
}
