// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// Windows containers are a special beast in Docker; you can't use docker cp
// to move files between the container and host.

// This communicator works around that limitation by reusing all possible
// methods and fields of the normal Docker Communicator, but we overwrite the
// Upload, Download, and UploadDir methods to utilize a mounted directory and
// native powershell commands rather than relying on docker cp.

type WindowsContainerCommunicator struct {
	Communicator
}

func powerShellSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// Creates the parent and all other preceding directories in the container for the given destination.
// Does not recreate paths that already exist.
func (c *WindowsContainerCommunicator) ensureContainerParentDir(ctx context.Context, destination string) error {
	cmd := &packersdk.RemoteCmd{
		Command: strings.Join([]string{
			fmt.Sprintf("Split-Path -Parent %s", powerShellSingleQuote(destination)),
			"Where-Object { $_ -and -not (Test-Path -LiteralPath $_) }",
			"ForEach-Object { New-Item -ItemType Directory -Force -Path $_ | Out-Null }",
		}, " | "),
	}
	if err := c.Start(ctx, cmd); err != nil {
		return err
	}

	cmd.Wait()
	if cmd.ExitStatus() != 0 {
		return fmt.Errorf("Upload failed to create parent directory for %s: %d", destination, cmd.ExitStatus())
	}

	return nil
}

// Upload uses docker exec to copy the file from the host to the container
func (c *WindowsContainerCommunicator) Upload(dst string, src io.Reader, fi *os.FileInfo) error {
	// Create a temporary file to store the upload
	tempfile, err := ioutil.TempFile(c.HostDir, "upload")
	if err != nil {
		return err
	}
	defer os.Remove(tempfile.Name())

	// Copy the contents to the temporary file
	_, err = io.Copy(tempfile, src)
	if err != nil {
		return err
	}
	if fi != nil {
		//nolint:errcheck
		tempfile.Chmod((*fi).Mode())
	}
	tempfile.Close()

	// Before copying the file into place, we need to make sure that the parent folders
	// exists if windows_create_parent_dirs was specified in the plugin config.
	// See also https://github.com/hashicorp/packer-plugin-docker/issues/208
	ctx := context.TODO()
	if c.Config.WindowsCreateParentDirs {
		if err := c.ensureContainerParentDir(ctx, dst); err != nil {
			return err
		}
	}

	// Copy the file into place by copying the temporary file we put
	// into the shared folder into the proper location in the container
	cmd := &packersdk.RemoteCmd{
		Command: fmt.Sprintf("Copy-Item -Path %s/%s -Destination %s", c.ContainerDir,
			filepath.Base(tempfile.Name()), dst),
	}
	if err := c.Start(ctx, cmd); err != nil {
		return err
	}

	// Wait for the copy to complete
	cmd.Wait()
	if cmd.ExitStatus() != 0 {
		return fmt.Errorf("Upload failed with non-zero exit status: %d", cmd.ExitStatus())
	}

	return nil
}

func (c *WindowsContainerCommunicator) UploadDir(dst string, src string, exclude []string) error {
	// Create the temporary directory that will store the contents of "src"
	// for copying into the container.
	td, err := ioutil.TempDir(c.HostDir, "dirupload")
	if err != nil {
		return err
	}
	defer os.RemoveAll(td)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relpath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		hostpath := filepath.Join(td, relpath)

		// If it is a directory, just create it
		if info.IsDir() {
			return os.MkdirAll(hostpath, info.Mode())
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			dest, err := os.Readlink(path)

			if err != nil {
				return err
			}

			return os.Symlink(dest, hostpath)
		}

		// It is a file, copy it over, including mode.
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(hostpath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return nil
	}

	// Copy the entire directory tree to the temporary directory
	if err := filepath.Walk(src, walkFn); err != nil {
		return err
	}

	// Determine the destination directory
	containerSrc := filepath.Join(c.ContainerDir, filepath.Base(td))
	containerDst := dst
	if src[len(src)-1] != '/' {
		containerDst = filepath.Join(dst, filepath.Base(src))
	}

	// Before copying the files into place, we need to make sure that the parent folders
	// exists if windows_create_parent_dirs was specified in the plugin config.
	// See also https://github.com/hashicorp/packer-plugin-docker/issues/208
	ctx := context.TODO()
	if c.Config.WindowsCreateParentDirs {
		if err := c.ensureContainerParentDir(ctx, containerDst); err != nil {
			return err
		}
	}

	// Make the directory, then copy into it
	cmd := &packersdk.RemoteCmd{
		Command: fmt.Sprintf("Copy-Item %s -Destination %s -Recurse",
			containerSrc, containerDst),
	}
	if err := c.Start(ctx, cmd); err != nil {
		return err
	}

	// Wait for the copy to complete
	cmd.Wait()
	if cmd.ExitStatus() != 0 {
		return fmt.Errorf("Upload failed with non-zero exit status: %d", cmd.ExitStatus())
	}

	return nil
}

// Download pulls a file out of a container using `docker cp`. We have a source
// path and want to write to an io.Writer
func (c *WindowsContainerCommunicator) Download(src string, dst io.Writer) error {
	log.Printf("Downloading file from container: %s:%s", c.ContainerID, src)
	// Copy file onto temp file on mounted volume inside container
	var stdout, stderr bytes.Buffer
	cmd := &packersdk.RemoteCmd{
		Command: fmt.Sprintf("Copy-Item -Path %s -Destination %s/%s", src, c.ContainerDir,
			filepath.Base(src)),
		Stdout: &stdout,
		Stderr: &stderr,
	}
	ctx := context.TODO()
	if err := c.Start(ctx, cmd); err != nil {
		return err
	}

	// Wait for the copy to complete
	cmd.Wait()

	if cmd.ExitStatus() != 0 {
		return fmt.Errorf("Failed to copy file to shared drive: %s, %s, %d", stderr.String(), stdout.String(), cmd.ExitStatus())
	}

	// Read that copied file into a new file opened on host machine
	fsrc, err := os.Open(filepath.Join(c.HostDir, filepath.Base(src)))
	if err != nil {
		return err
	}
	defer fsrc.Close()
	defer os.Remove(fsrc.Name())

	_, err = io.Copy(dst, fsrc)
	if err != nil {
		return err
	}

	return nil
}
