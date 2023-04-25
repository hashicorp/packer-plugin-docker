// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import "testing"

func TestDockerDriver_impl(t *testing.T) {
	var _ Driver = new(DockerDriver)
}
