package docker

import (
	"context"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"
)

func TestStepSetGeneratedData_Run(t *testing.T) {
	state := testState(t)
	step := &StepDefaultGeneratedData{
		GeneratedData: &packerbuilderdata.GeneratedData{State: state},
	}

	if action := step.Run(context.TODO(), state); action != multistep.ActionContinue {
		t.Fatalf("Should not halt")
	}

	genData := state.Get("generated_data").(map[string]interface{})
	imgSha256 := genData["ImageSha256"].(string)
	if imgSha256 != "ERR_IMAGE_SHA256_NOT_FOUND" {
		t.Fatalf("Expected ImageSha256 to be ERR_IMAGE_SHA256_NOT_FOUND but was %s", imgSha256)
	}

	sourceDigest := genData["SourceImageDigest"].(string)
	if sourceDigest != "ERR_SOURCE_IMAGE_DIGEST_NOT_FOUND" {
		t.Fatalf("Expected SourceImageDigest to be ERR_SOURCE_IMAGE_DIGEST_NOT_FOUND but was %s", sourceDigest)
	}

	sourceSha256 := genData["SourceImageSha256"].(string)
	if sourceSha256 != "ERR_SOURCE_IMAGE_Sha256_NOT_FOUND" {
		t.Fatalf("Expected SourceImageSha256 to be ERR_SOURCE_IMAGE_Sha256_NOT_FOUND but was %s", sourceSha256)
	}
}
