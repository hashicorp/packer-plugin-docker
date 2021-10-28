package docker

import (
	"errors"
	"strings"
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	registryimage "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
	"github.com/mitchellh/mapstructure"
)

func TestImportArtifact_impl(t *testing.T) {
	var _ packersdk.Artifact = new(ImportArtifact)
}

func TestImportArtifactBuilderId(t *testing.T) {
	a := &ImportArtifact{BuilderIdValue: "foo"}
	if a.BuilderId() != "foo" {
		t.Fatalf("bad: %#v", a.BuilderId())
	}
}

func TestImportArtifactFiles(t *testing.T) {
	a := &ImportArtifact{}
	if a.Files() != nil {
		t.Fatalf("bad: %#v", a.Files())
	}
}

func TestImportArtifactId(t *testing.T) {
	a := &ImportArtifact{IdValue: "foo"}
	if a.Id() != "foo" {
		t.Fatalf("bad: %#v", a.Id())
	}
}

func TestImportArtifactDestroy(t *testing.T) {
	d := new(MockDriver)
	a := &ImportArtifact{
		Driver:  d,
		IdValue: "foo",
	}

	// No error
	if err := a.Destroy(); err != nil {
		t.Fatalf("err: %s", err)
	}
	if !d.DeleteImageCalled {
		t.Fatal("delete image should be called")
	}
	if d.DeleteImageId != "foo" {
		t.Fatalf("bad: %#v", d.DeleteImageId)
	}

	// With an error
	d.DeleteImageErr = errors.New("foo")
	if err := a.Destroy(); err != d.DeleteImageErr {
		t.Fatalf("err: %#v", err)
	}
}

func TestArtifactState_RegistryImageMetadataNoGenData(t *testing.T) {

	artifact := &ImportArtifact{
		Driver:  new(MockDriver),
		IdValue: "docker-image",
		StateData: map[string]interface{}{
			"docker_tags": []string{"tag1", "tag2"},
		},
	}
	// Valid state
	result := artifact.State(registryimage.ArtifactStateURI)
	if result == nil {
		t.Fatalf("Bad: HCP Packer registry image data was nil")
	}

	var image registryimage.Image
	err := mapstructure.Decode(result, &image)
	if err != nil {
		t.Errorf("Bad: unexpected error when trying to decode state into registryimage.Image %v", err)
	}

	if image.ImageID != artifact.IdValue {
		t.Errorf("Bad: unexpected value for ImageID %q, expected %q", image.ImageID, artifact.IdValue)
	}

	if image.Labels["tags"] != strings.Join(artifact.loadTags(), ",") {
		t.Errorf("Bad: unexpected value for tags but got %v", image.Labels["tags"])
	}

	if image.ProviderRegion != "docker" {
		t.Errorf("Bad: unexpected value for ImageID %q, expected docker", image.ProviderRegion)
	}

}

func TestArtifactState_RegistryImageMetadataWithGenData(t *testing.T) {

	artifact := &ImportArtifact{
		Driver:  new(MockDriver),
		IdValue: "docker-image",
		StateData: map[string]interface{}{
			"docker_tags": []string{"tag1", "tag2"},
		},
	}
	genData := map[string]interface{}{
		"SourceImageSha256": "sha256",
		"SourceImageDigest": "digest",
		"ImageSha256":       "ImageSha256String",
		"PackerArtifactID":  artifact.Id(),
	}
	artifact.StateData["generated_data"] = genData

	// Valid state
	result := artifact.State(registryimage.ArtifactStateURI)
	if result == nil {
		t.Fatalf("Bad: HCP Packer registry image data was nil")
	}

	var image registryimage.Image
	err := mapstructure.Decode(result, &image)
	if err != nil {
		t.Errorf("Bad: unexpected error when trying to decode state into registryimage.Image %v", err)
	}

	if image.ImageID != "ImageSha256String" {
		t.Errorf("Bad: unexpected value for ImageID %q", image.ImageID)
	}

	if image.Labels["tags"] != strings.Join(artifact.loadTags(), ",") {
		t.Errorf("Bad: unexpected value for tags but got %v", image.Labels["tags"])
	}

	if image.ProviderRegion != "docker" {
		t.Errorf("Bad: unexpected value for ImageID %q, expected docker", image.ProviderRegion)
	}

	for k, v := range genData {
		if k == "SourceImageSha256" {
			continue
		}

		if image.Labels[k] != v {
			t.Errorf("Bad: unexpected value for label %q, expected %q but got %q", k, genData[k].(string), image.Labels[k])
		}
	}

}
