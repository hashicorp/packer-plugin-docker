// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/template/config"
)

func TestBuildConfigIsDefault(t *testing.T) {
	tests := []struct {
		name          string
		inputConfig   DockerfileBootstrapConfig
		expectDefault bool
	}{
		{
			"empty config - should be default",
			DockerfileBootstrapConfig{},
			true,
		},
		{
			"dockerfile set - should not be default",
			DockerfileBootstrapConfig{
				DockerfilePath: "test",
			},
			false,
		},
		{
			"pull explicitly set to false - should not be default",
			DockerfileBootstrapConfig{
				Pull: config.TriFalse,
			},
			false,
		},
		{
			"build dir set - should not be default",
			DockerfileBootstrapConfig{
				BuildDir: "dir",
			},
			false,
		},
		{
			"empty argument map - should be default",
			DockerfileBootstrapConfig{
				Arguments: map[string]string{},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isDefault := tt.inputConfig.IsDefault()
			if isDefault != tt.expectDefault {
				t.Errorf("expected isDefault %t is different from reported %t",
					tt.expectDefault, isDefault)
			}
		})
	}
}
