// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

// TestUploadDownload verifies that basic upload / download functionality works
func TestUploadDownload(t *testing.T) {
	if os.Getenv("PACKER_ACC") == "" {
		t.Skip("This test is only run with PACKER_ACC=1")
	}

	testcases := []string{
		"upload_download.json",
		"upload_download.pkr.hcl",
	}

	for _, tc := range testcases {
		templatePath := filepath.Join("test-fixtures", tc)
		bytes, err := ioutil.ReadFile(templatePath)
		if err != nil {
			t.Fatalf("failed to load template file %s", templatePath)
		}
		configString := string(bytes)

		// this should be a precheck
		cmd := exec.Command("docker", "-v")
		err = cmd.Run()
		if err != nil {
			t.Error("docker command not found; please make sure docker is installed")
		}

		acctest.TestPlugin(t, &acctest.PluginTestCase{
			Name:     "Upload_Download",
			Template: configString,
			Type:     "docker",
			Check: func(buildCommand *exec.Cmd, logfile string) error {
				// Verify that the thing we downloaded is the same thing we sent up.
				// Complain loudly if it isn't.
				inputFile, err := ioutil.ReadFile("test-fixtures/onecakes/strawberry")
				if err != nil {
					return fmt.Errorf("Unable to read input file: %s", err)
				}
				outputFile, err := ioutil.ReadFile("my-strawberry-cake")
				if err != nil {
					return fmt.Errorf("Unable to read output file: %s", err)
				}
				if sha256.Sum256(inputFile) != sha256.Sum256(outputFile) {
					return fmt.Errorf("Input and output files do not match\n"+
						"Input:\n%s\nOutput:\n%s\n", inputFile, outputFile)
				}
				return nil
			},
			Teardown: func() error {
				// Cleanup. Honestly I don't know why you would want to get rid
				// of my strawberry cake. It's so tasty! Do you not like cake? Are you a
				// cake-hater? Or are you keeping all the cake all for yourself? So selfish!
				os.Remove("my-strawberry-cake")
				return nil
			},
		})
	}
}

// TestLargeDownload verifies that files are the appropriate size after being
// downloaded. This is to identify and fix the race condition in #2793. You may
// need to use github.com/cbednarski/rerun to verify since this problem occurs
// only intermittently.
func TestLargeDownload(t *testing.T) {
	if os.Getenv("PACKER_ACC") == "" {
		t.Skip("This test is only run with PACKER_ACC=1")
	}

	testcases := []string{
		"large_download.json",
		"large_download.pkr.hcl",
	}
	for _, tc := range testcases {
		templatePath := filepath.Join("test-fixtures", tc)
		bytes, err := ioutil.ReadFile(templatePath)
		if err != nil {
			t.Fatalf("failed to load template file %s", templatePath)
		}
		configString := string(bytes)

		// this should be a precheck
		cmd := exec.Command("docker", "-v")
		err = cmd.Run()
		if err != nil {
			t.Error("docker command not found; please make sure docker is installed")
		}

		acctest.TestPlugin(t, &acctest.PluginTestCase{
			Name:     "Large_Download",
			Type:     "docker",
			Template: configString,
			Check: func(buildCommand *exec.Cmd, logfile string) error {
				// Verify that the things we downloaded are the right size. Complain loudly
				// if they are not.
				//
				// cupcake should be 2097152 bytes
				// bigcake should be 104857600 bytes
				cupcake, err := os.Stat("cupcake")
				if err != nil {
					t.Fatalf("Unable to stat cupcake file: %s", err)
				}
				cupcakeExpected := int64(2097152)
				if cupcake.Size() != cupcakeExpected {
					t.Errorf("Expected cupcake to be %d bytes; found %d", cupcakeExpected, cupcake.Size())
				}

				bigcake, err := os.Stat("bigcake")
				if err != nil {
					t.Fatalf("Unable to stat bigcake file: %s", err)
				}
				bigcakeExpected := int64(104857600)
				if bigcake.Size() != bigcakeExpected {
					t.Errorf("Expected bigcake to be %d bytes; found %d", bigcakeExpected, bigcake.Size())
				}

				// TODO if we can, calculate a sha inside the container and compare to the
				// one we get after we pull it down. We will probably have to parse the log
				// or ui output to do this because we use /dev/urandom to create the file.

				// if sha256.Sum256(inputFile) != sha256.Sum256(outputFile) {
				//	t.Fatalf("Input and output files do not match\n"+
				//		"Input:\n%s\nOutput:\n%s\n", inputFile, outputFile)
				// }
				return nil
			},
			Teardown: func() error {
				os.Remove("cupcake")
				os.Remove("bigcake")
				return nil
			},
		})
	}

}

// TestFixUploadOwner verifies that owner of uploaded files is the user the  container is running as.
func TestFixUploadOwner(t *testing.T) {
	if os.Getenv("PACKER_ACC") == "" {
		t.Skip("This test is only run with PACKER_ACC=1")
	}

	cmd := exec.Command("docker", "-v")
	err := cmd.Run()
	if err != nil {
		t.Error("docker command not found; please make sure docker is installed")
	}

	testcases := []string{
		"fix_upload_owner.json",
		"fix_upload_owner.pkr.hcl",
	}
	for _, tc := range testcases {
		templatePath := filepath.Join("test-fixtures", tc)
		bytes, err := ioutil.ReadFile(templatePath)
		if err != nil {
			t.Fatalf("failed to load template file %s", templatePath)
		}
		configString := string(bytes)

		acctest.TestPlugin(t, &acctest.PluginTestCase{
			Template: configString,
		})
	}
}
