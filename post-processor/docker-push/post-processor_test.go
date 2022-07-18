package dockerpush

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/packer-plugin-docker/builder/docker"
	dockerimport "github.com/hashicorp/packer-plugin-docker/post-processor/docker-import"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

func testUi() *packersdk.BasicUi {
	return &packersdk.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	}
}

// This reads the output from the bytes.Buffer in our test UI
// and then resets the buffer.
func readWriter(ui *packersdk.BasicUi) (resultString string) {
	buffer := ui.Writer.(*bytes.Buffer)
	resultString = buffer.String()
	buffer.Reset()
	return
}

func TestPostProcessor_ImplementsPostProcessor(t *testing.T) {
	var _ packersdk.PostProcessor = new(PostProcessor)
}

func TestGetEcrType(t *testing.T) {
	awsConfig := docker.AwsAccessConfig{}

	typ, _ := awsConfig.GetEcrType("https://public.ecr.aws/j9y7g6y8/dev_hc_pkr_dkr_test_1")
	if typ != docker.Public {
		msg := fmt.Sprintf("ECR type should be %v", docker.Public)
		t.Fatal(msg)
	}

	typ, _ = awsConfig.GetEcrType("https://public.ecr.aws/j9y7g6y8/dev_hc_pkr_dkr_test_1")
	if typ != docker.Public {
		msg := fmt.Sprintf("ECR type should be %v", docker.Public)
		t.Fatal(msg)
	}

	typ, _ = awsConfig.GetEcrType("https://12345.dkr.ecr.us-east-1.amazonaws.com/private_dev_hc_pkr_dkr_test_1")
	if typ != docker.Private {
		msg := fmt.Sprintf("ECR type should be %v", docker.Private)
		t.Fatal(msg)
	}

	typ, _ = awsConfig.GetEcrType("google.com")
	if typ != docker.Invalid {
		msg := fmt.Sprintf("ECR type should be %v", docker.Invalid)
		t.Fatal(msg)
	}
}

func TestPostProcessor_PostProcess(t *testing.T) {
	driver := &docker.MockDriver{}
	p := &PostProcessor{Driver: driver}
	artifact := &packersdk.MockArtifact{
		BuilderIdValue: dockerimport.BuilderId,
		IdValue:        "foo/bar",
	}

	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi(), artifact)
	if _, ok := result.(packersdk.Artifact); !ok {
		t.Fatal("should be instance of Artifact")
	}
	if !keep {
		t.Fatal("should keep")
	}
	if forceOverride {
		t.Fatal("Should default to keep, but not override user wishes")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !driver.PushCalled {
		t.Fatal("should call push")
	}
	if driver.PushName != "foo/bar" {
		t.Fatal("bad name")
	}
	if result.Id() != "foo/bar" {
		t.Fatal("bad image id")
	}
}

func TestPostProcessor_PostProcess_portInName(t *testing.T) {
	driver := &docker.MockDriver{}
	p := &PostProcessor{Driver: driver}
	artifact := &packersdk.MockArtifact{
		BuilderIdValue: dockerimport.BuilderId,
		IdValue:        "localhost:5000/foo/bar",
	}

	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi(), artifact)
	if _, ok := result.(packersdk.Artifact); !ok {
		t.Fatal("should be instance of Artifact")
	}
	if !keep {
		t.Fatal("should keep")
	}
	if forceOverride {
		t.Fatal("Should default to keep, but not override user wishes")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !driver.PushCalled {
		t.Fatal("should call push")
	}
	if driver.PushName != "localhost:5000/foo/bar" {
		t.Fatal("bad name")
	}
	if result.Id() != "localhost:5000/foo/bar" {
		t.Fatal("bad image id")
	}
}

func TestPostProcessor_PostProcess_tags(t *testing.T) {
	driver := &docker.MockDriver{}
	p := &PostProcessor{Driver: driver}
	artifact := &packersdk.MockArtifact{
		BuilderIdValue: dockerimport.BuilderId,
		IdValue:        "hashicorp/ubuntu:precise",
	}

	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi(), artifact)
	if _, ok := result.(packersdk.Artifact); !ok {
		t.Fatal("should be instance of Artifact")
	}
	if !keep {
		t.Fatal("should keep")
	}
	if forceOverride {
		t.Fatal("Should default to keep, but not override user wishes")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !driver.PushCalled {
		t.Fatal("should call push")
	}
	if driver.PushName != "hashicorp/ubuntu:precise" {
		t.Fatalf("bad name: %s", driver.PushName)
	}
	if result.Id() != "hashicorp/ubuntu:precise" {
		t.Fatal("bad image id")
	}
}

func TestPostProcessor_PostProcess_digestWarning(t *testing.T) {
	driver := &docker.MockDriver{}
	p := &PostProcessor{Driver: driver}
	artifact := &packersdk.MockArtifact{
		BuilderIdValue: dockerimport.BuilderId,
		IdValue:        "hashicorp/ubuntu:precise",
	}

	driver.DigestErr = fmt.Errorf("I'm a generic digest error! The Packer Docker Plugin should handle me as a warning")

	testUi := testUi()
	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi, artifact)
	resultString := readWriter(testUi)
	if _, ok := result.(packersdk.Artifact); !ok {
		t.Fatal("should be instance of Artifact")
	}
	if !keep {
		t.Fatal("should keep")
	}
	if forceOverride {
		t.Fatal("Should default to keep, but not override user wishes")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !driver.PushCalled {
		t.Fatal("should call push")
	}
	// Check for warning text
	if !strings.Contains(resultString, "Unable to determine digest for source image, ignoring it for now") {
		t.Fatal(resultString)
	}
	// Should still succeed after digest warning
	if driver.PushName != "hashicorp/ubuntu:precise" {
		t.Fatalf("bad name: %s", driver.PushName)
	}
	if result.Id() != "hashicorp/ubuntu:precise" {
		t.Fatal("bad image id")
	}
}
