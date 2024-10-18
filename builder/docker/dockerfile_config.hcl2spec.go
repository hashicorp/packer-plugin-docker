// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package docker

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatDockerfileBootstrapConfig is an auto-generated flat version of DockerfileBootstrapConfig.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatDockerfileBootstrapConfig struct {
	DockerfilePath *string `mapstructure:"path" required:"true" cty:"path" hcl:"path"`
	BuildDir       *string `mapstructure:"build_dir" cty:"build_dir" hcl:"build_dir"`
	Pull           *bool   `mapstructure:"pull" cty:"pull" hcl:"pull"`
	Compress       *bool   `mapstructure:"compress" cty:"compress" hcl:"compress"`
}

// FlatMapstructure returns a new FlatDockerfileBootstrapConfig.
// FlatDockerfileBootstrapConfig is an auto-generated flat version of DockerfileBootstrapConfig.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*DockerfileBootstrapConfig) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatDockerfileBootstrapConfig)
}

// HCL2Spec returns the hcl spec of a DockerfileBootstrapConfig.
// This spec is used by HCL to read the fields of DockerfileBootstrapConfig.
// The decoded values from this spec will then be applied to a FlatDockerfileBootstrapConfig.
func (*FlatDockerfileBootstrapConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"path":      &hcldec.AttrSpec{Name: "path", Type: cty.String, Required: false},
		"build_dir": &hcldec.AttrSpec{Name: "build_dir", Type: cty.String, Required: false},
		"pull":      &hcldec.AttrSpec{Name: "pull", Type: cty.Bool, Required: false},
		"compress":  &hcldec.AttrSpec{Name: "compress", Type: cty.Bool, Required: false},
	}
	return s
}
