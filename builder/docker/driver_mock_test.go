// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import "testing"

func TestMockDriver_impl(t *testing.T) {
	var _ Driver = new(MockDriver)
}
