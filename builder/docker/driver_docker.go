// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/go-version"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type DockerDriver struct {
	Ui  packersdk.Ui
	Ctx *interpolate.Context

	// The directory Docker should use to store its client configuration.
	// Provides an isolated client configuration to each Docker operation to
	// prevent race conditions.
	ConfigDir string
	// The executable to run commands with.
	Executable string

	l sync.Mutex
}

func (d *DockerDriver) Build(args []string) (string, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	imageIdFile, err := os.CreateTemp("", "")
	if err != nil {
		return "", fmt.Errorf("failed to create image ID file: %s", err)
	}
	imageIdFilePath := imageIdFile.Name()
	imageIdFile.Close()

	cmd := exec.Command(d.Executable, "build")
	cmd.Args = append(cmd.Args, "--iidfile", imageIdFilePath)
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s build failed: %s; stdout: %s; stderr: %s", d.Executable, err, stdout.String(), stderr.String())
	}

	imageId, err := os.ReadFile(imageIdFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image ID from file %q: %s", imageIdFilePath, err)
	}

	return strings.TrimSpace(string(imageId)), nil
}

func (d *DockerDriver) DeleteImage(id string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(d.Executable, "rmi", id)
	cmd.Stderr = &stderr

	log.Printf("Deleting image: %s", id)
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error deleting image: %s\nStderr: %s",
			err, stderr.String())
		return err
	}

	return nil
}

func (d *DockerDriver) Commit(id string, author string, changes []string, message string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	args := []string{"commit"}
	if author != "" {
		args = append(args, "--author", author)
	}
	for _, change := range changes {
		args = append(args, "--change", change)
	}
	if message != "" {
		args = append(args, "--message", message)
	}
	args = append(args, id)

	log.Printf("Committing container with args: %v", args)
	cmd := exec.Command(d.Executable, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error committing container: %s\nStderr: %s",
			err, stderr.String())
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (d *DockerDriver) Export(id string, dst io.Writer) error {
	var stderr bytes.Buffer
	cmd := exec.Command(d.Executable, "export", id)
	cmd.Stdout = dst
	cmd.Stderr = &stderr

	log.Printf("Exporting container: %s", id)
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error exporting: %s\nStderr: %s",
			err, stderr.String())
		return err
	}

	return nil
}

func (d *DockerDriver) Import(path string, changes []string, repo string, platform string) (string, error) {
	var stdout, stderr bytes.Buffer

	args := []string{"import"}

	for _, change := range changes {
		args = append(args, "--change", change)
	}

	if platform != "" {
		args = append(args, "--platform", platform)
	}

	args = append(args, "-")
	args = append(args, repo)

	cmd := exec.Command(d.Executable, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return "", err
	}

	// There should be only one artifact of the Docker builder
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	log.Printf("Importing tarball with args: %v", args)

	if err := cmd.Start(); err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		//nolint
		io.Copy(stdin, file)
	}()

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("Error importing container: %s\n\nStderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (d *DockerDriver) IPAddress(id string) (string, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command(
		d.Executable,
		"inspect",
		"--format",
		"{{ .NetworkSettings.IPAddress }}",
		id)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error: %s\n\nStderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// Sha256 retrieves the image Id using Docker inspect.
func (d *DockerDriver) Sha256(id string) (string, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command(
		d.Executable,
		"inspect",
		"--format",
		"{{ .Id }}",
		id)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error: %s\n\nStderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// Digest retrieves the digest of the image using Docker inspect.
// Format for the digest is: <repo>@sha256:<shasum>
// For example:
// ubuntu@sha256:454054f5bbd571b088db25b662099c6c7b3f0cb78536a2077d54adc48f00cd68
// This can be considered a source of truth for pointing to a specific image
// at a specific point in time.
func (d *DockerDriver) Digest(id string) (string, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command(
		d.Executable,
		"inspect",
		"--format",
		"{{ ( index .RepoDigests 0 ) }}",
		id)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error: %s\n\nStderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (d *DockerDriver) Cmd(id string) (string, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command(
		d.Executable,
		"inspect",
		"--format",
		"{{if .Config.Cmd}} {{json .Config.Cmd}} {{else}} [\"\"] {{end}}",
		id)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error: %s\n\nStderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (d *DockerDriver) Entrypoint(id string) (string, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command(
		d.Executable,
		"inspect",
		"--format",
		"{{if .Config.Entrypoint}} {{json .Config.Entrypoint}} {{else}} [\"\"] {{end}}",
		id)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error: %s\n\nStderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (d *DockerDriver) Login(repo, user, pass string) error {
	d.l.Lock()

	version_running, err := d.Version()
	if err != nil {
		d.l.Unlock()
		return err
	}

	// Version 17.07.0 of Docker adds support for the new
	// `--password-stdin` option which can be used to offer
	// password via the standard input, rather than passing
	// the password and/or token using a command line switch.
	constraint, err := version.NewConstraint(">= 17.07.0")
	if err != nil {
		d.l.Unlock()
		return err
	}

	cmd := d.newCommandWithConfig("login")

	if user != "" {
		cmd.Args = append(cmd.Args, "-u", user)
	}

	if pass != "" {
		if constraint.Check(version_running) {
			cmd.Args = append(cmd.Args, "--password-stdin")

			stdin, err := cmd.StdinPipe()
			if err != nil {
				d.l.Unlock()
				return err
			}
			_, err = io.WriteString(stdin, pass)
			if err != nil {
				return err
			}
			stdin.Close()
		} else {
			cmd.Args = append(cmd.Args, "-p", pass)
		}
	}

	if repo != "" {
		cmd.Args = append(cmd.Args, repo)
	}

	err = runAndStream(cmd, d.Ui)
	if err != nil {
		d.l.Unlock()
		return err
	}

	return nil
}

func (d *DockerDriver) Logout(repo string) error {
	cmd := d.newCommandWithConfig("logout")

	if repo != "" {
		cmd.Args = append(cmd.Args, repo)
	}

	err := runAndStream(cmd, d.Ui)
	d.l.Unlock()
	return err
}

func (d *DockerDriver) Pull(image string, platform string) error {
	cmd := d.newCommandWithConfig("pull", image)

	if platform != "" {
		cmd.Args = append(cmd.Args, "--platform", platform)
	}

	return runAndStream(cmd, d.Ui)
}

func (d *DockerDriver) Push(name string, platform string) error {
	cmd := d.newCommandWithConfig("push", name)

	if platform != "" {
		cmd.Args = append(cmd.Args, "--platform", platform)
	}

	return runAndStream(cmd, d.Ui)
}

func (d *DockerDriver) SaveImage(id string, dst io.Writer) error {
	var stderr bytes.Buffer
	cmd := exec.Command(d.Executable, "save", id)
	cmd.Stdout = dst
	cmd.Stderr = &stderr

	log.Printf("Exporting image: %s", id)
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error exporting: %s\nStderr: %s",
			err, stderr.String())
		return err
	}

	return nil
}

func (d *DockerDriver) StartContainer(config *ContainerConfig) (string, error) {
	// Build up the template data
	var tplData startContainerTemplate
	tplData.Image = config.Image
	ictx := *d.Ctx
	ictx.Data = &tplData

	// Args that we're going to pass to Docker
	args := []string{"run"}
	for _, v := range config.Device {
		args = append(args, "--device", v)
	}
	for _, v := range config.CapAdd {
		args = append(args, "--cap-add", v)
	}
	for _, v := range config.CapDrop {
		args = append(args, "--cap-drop", v)
	}
	if config.Privileged {
		args = append(args, "--privileged")
	}
	if config.Runtime != "" {
		args = append(args, "--runtime", config.Runtime)
	}
	if config.Platform != "" {
		args = append(args, "--platform", config.Platform)
	}
	for _, v := range config.TmpFs {
		args = append(args, "--tmpfs", v)
	}
	for host, guest := range config.Volumes {
		if strings.HasPrefix(host, "~/") {
			homedir, _ := os.UserHomeDir()
			host = filepath.Join(homedir, host[2:])
		}
		args = append(args, "-v", fmt.Sprintf("%s:%s", host, guest))
	}
	for _, v := range config.RunCommand {
		v, err := interpolate.Render(v, &ictx)
		if err != nil {
			return "", err
		}

		args = append(args, v)
	}
	d.Ui.Message(fmt.Sprintf(
		"Run command: %s %s", d.Executable, strings.Join(args, " ")))

	// Start the container
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(d.Executable, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Printf("Starting container with args: %v", args)
	if err := cmd.Start(); err != nil {
		return "", err
	}

	log.Println("Waiting for container to finish starting")
	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf("Docker exited with a non-zero exit status.\nStderr: %s",
				stderr.String())
		}

		return "", err
	}

	// Capture the container ID, which is alone on stdout
	return strings.TrimSpace(stdout.String()), nil
}

func (d *DockerDriver) StopContainer(id string) error {
	if err := exec.Command(d.Executable, "stop", id).Run(); err != nil {
		return err
	}
	return nil
}

func (d *DockerDriver) KillContainer(id string) error {
	if err := exec.Command(d.Executable, "kill", id).Run(); err != nil {
		return err
	}

	return exec.Command(d.Executable, "rm", id).Run()
}

func (d *DockerDriver) TagImage(id string, repo string, force bool) error {
	args := []string{"tag"}

	// detect running docker version before tagging
	// flag `force` for docker tagging was removed after Docker 1.12.0
	// to keep its backward compatibility, we are not going to remove `force`
	// option, but to ignore it when Docker version >= 1.12.0
	//
	// for more detail, please refer to the following links:
	// - https://docs.docker.com/engine/deprecated/#/f-flag-on-docker-tag
	// - https://github.com/docker/docker/pull/23090
	version_running, err := d.Version()
	if err != nil {
		return err
	}

	version_deprecated, err := version.NewVersion("1.12.0")
	if err != nil {
		// should never reach this line
		return err
	}

	if force {
		if version_running.LessThan(version_deprecated) {
			args = append(args, "-f")
		} else {
			// do nothing if Docker version >= 1.12.0
			log.Printf("[WARN] option: \"force\" will be ignored here")
			log.Printf("since it was removed after Docker 1.12.0 released")
		}
	}
	args = append(args, id, repo)

	var stderr bytes.Buffer
	cmd := exec.Command(d.Executable, args...)
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error tagging image: %s\nStderr: %s",
			err, stderr.String())
		return err
	}

	return nil
}

func (d *DockerDriver) Verify() error {
	if _, err := exec.LookPath(d.Executable); err != nil {
		return err
	}

	return nil
}

func (d *DockerDriver) Version() (*version.Version, error) {
	output, err := exec.Command(d.Executable, "-v").Output()
	if err != nil {
		return nil, err
	}

	match := regexp.MustCompile(version.VersionRegexpRaw).FindSubmatch(output)
	if match == nil {
		return nil, fmt.Errorf("unknown version: %s", output)
	}

	log.Printf("version matches: %s", match)

	return version.NewVersion(string(match[0]))
}

func (d *DockerDriver) newCommandWithConfig(args ...string) *exec.Cmd {
	cmd := exec.Command(d.Executable)

	if d.ConfigDir != "" {
		cmd.Args = append(cmd.Args, "--config", d.ConfigDir)
	}

	cmd.Args = append(cmd.Args, args...)

	return cmd
}
