// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestExportArtifact_impl(t *testing.T) {
	var _ packersdk.Artifact = new(ExportArtifact)
}
