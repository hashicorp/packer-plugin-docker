// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dockersave

import (
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestPostProcessor_ImplementsPostProcessor(t *testing.T) {
	var _ packersdk.PostProcessor = new(PostProcessor)
}
