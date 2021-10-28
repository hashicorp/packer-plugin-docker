package docker

import (
	"fmt"
	"log"
	"strings"

	registryimage "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
)

// ImportArtifact is an Artifact implementation for when a container is
// exported from docker into a single flat file.
type ImportArtifact struct {
	BuilderIdValue string
	Driver         Driver
	IdValue        string

	// StateData should store data such as GeneratedData
	// to be shared with post-processors
	StateData map[string]interface{}
}

func (a *ImportArtifact) BuilderId() string {
	return a.BuilderIdValue
}

func (*ImportArtifact) Files() []string {
	return nil
}

func (a *ImportArtifact) Id() string {
	return a.IdValue
}

func (a *ImportArtifact) String() string {
	tags := a.loadTags()
	if len(tags) > 0 {
		return fmt.Sprintf("Imported Docker image: %s with tags %s",
			a.Id(), strings.Join(tags, " "))
	}
	return fmt.Sprintf("Imported Docker image: %s", a.Id())
}

func (a *ImportArtifact) State(name string) interface{} {
	if name == registryimage.ArtifactStateURI {
		return a.stateHCPPackerRegistryMetadata()
	}
	return a.StateData[name]
}

func (a *ImportArtifact) Destroy() error {
	return a.Driver.DeleteImage(a.Id())
}

func (a *ImportArtifact) loadTags() []string {
	var tags []string
	switch t := a.StateData["docker_tags"].(type) {
	case []string:
		tags = t
	case []interface{}:
		for _, name := range t {
			if n, ok := name.(string); ok {
				tags = append(tags, n)
			}
		}
	}
	return tags
}

// stateHCPPackerRegistryMetadata will write the metadata as an hcpRegistryImage
// Some notes on the Docker artifact. The image ID cannot be pulled using
// Docker Pull. However, it can be used as the source image to the docker
// builder if `pull` is set to false in the builder config. The digest can be
// pulled, but is not set by a `docker commit`, and is only set when
// `docker push` is called. For this reason, we are going to use the Id as the
// format for the source and stored output images, but also store the digest
// as a label so that users can retrieve the digest for use as a source image
// if they prefer.
func (a *ImportArtifact) stateHCPPackerRegistryMetadata() interface{} {
	labels := make(map[string]interface{})

	tags := a.loadTags()
	if len(tags) > 0 {
		labels["tags"] = strings.Join(tags, ",")
	}

	img, _ := registryimage.FromArtifact(a,
		registryimage.WithRegion("docker"),
		registryimage.WithProvider("docker"),
		registryimage.SetLabels(labels),
	)

	data, ok := a.StateData["generated_data"].(map[string]interface{})
	if !ok {
		log.Printf("No generated data exists in state. Artifact: %#v", a)
		return img
	}

	img.SourceImageID = data["SourceImageSha256"].(string)
	img.Labels["SourceImageDigest"] = data["SourceImageDigest"].(string)
	// This is the image's sha that we store as the image id. We store it
	// here as well becasue there is no guarantee this is the value stored
	// on the main artifact id value.
	img.Labels["ImageSha256"] = data["ImageSha256"].(string)
	// The docker tag and docker push post-processors store the repo:tag
	// combination here, whereas the docker builder stores the image's
	// sha256 id. We store this for posterity, but the image id needs to be
	// the sha256.
	img.Labels["PackerArtifactID"] = a.Id()
	// Overwrite ID with
	img.ImageID = data["ImageSha256"].(string)

	return img
}
