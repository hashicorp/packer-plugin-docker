Type: `docker`

The `docker` Packer builder builds [Docker](https://www.docker.io) images using
Docker. The builder starts a Docker container, runs provisioners within this
container, then exports the container for reuse or commits the image.

Packer builds Docker containers _without_ the use of
[Dockerfiles](https://docs.docker.com/engine/reference/builder/). By not using
`Dockerfiles`, Packer is able to provision containers with portable scripts or
configuration management systems that are not tied to Docker in any way. It
also has a simple mental model: you provision containers much the same way you
provision a normal virtualized or dedicated server. For more information, read
the section on [Dockerfiles](#dockerfiles).

The Docker builder must run on a machine that has Docker Engine installed.
Therefore the builder only works on machines that support Docker and _does not
support running on a Docker remote host_. You can learn about what [platforms
Docker supports and how to install onto
them](https://docs.docker.com/engine/installation/) in the Docker
documentation.

## Basic Example: Export

Below is a fully functioning example. It doesn't do anything useful, since no
provisioners are defined, but it will effectively repackage an image.

**HCL2**

```hcl
source "docker" "example" {
  image = "ubuntu"
  export_path = "image.tar"
}

**JSON**

```json
{
  "type": "docker",
  "image": "ubuntu",
  "export_path": "image.tar"
}
```

build {
  sources = ["source.docker.example"]
}
```

## Basic Example: Commit

Below is another example, the same as above but instead of exporting the
running container, this one commits the container to an image. The image can
then be more easily tagged, pushed, etc.

**HCL2**

```hcl
source "docker" "example" {
  image = "ubuntu"
  commit = true
}

build {
  sources = ["source.docker.example"]
}
```

**JSON**

```json
{
  "type": "docker",
  "image": "ubuntu",
  "export_path": "image.tar"
}
```

## Basic Example: Changes to Metadata

Below is an example using the changes argument of the builder. This feature
allows the source images metadata to be changed when committed back into the
Docker environment. It is derived from the `docker commit --change` command
line [option to
Docker](https://docs.docker.com/engine/reference/commandline/commit/).

Example uses of all of the options, assuming one is building an NGINX image
from ubuntu as an simple example:

**HCL2**

```hcl
source "docker" "example" {
    image = "ubuntu"
    commit = true
      changes = [
      "USER www-data",
      "WORKDIR /var/www",
      "ENV HOSTNAME www.example.com",
      "VOLUME /test1 /test2",
      "EXPOSE 80 443",
      "LABEL version=1.0",
      "ONBUILD RUN date",
      "CMD [\"nginx\", \"-g\", \"daemon off;\"]",
      "ENTRYPOINT /var/www/start.sh"
    ]
}
```

**JSON**

```json
{
  "type": "docker",
  "image": "ubuntu",
  "commit": true,
  "changes": [
    "USER www-data",
    "WORKDIR /var/www",
    "ENV HOSTNAME www.example.com",
    "VOLUME /test1 /test2",
    "EXPOSE 80 443",
    "LABEL version=1.0",
    "ONBUILD RUN date",
    "CMD [\"nginx\", \"-g\", \"daemon off;\"]",
    "ENTRYPOINT /var/www/start.sh"
  ]
}
```

Allowed metadata fields that can be changed are:

- CMD
  - String, supports both array (escaped) and string form
  - EX: `"CMD [\"nginx\", \"-g\", \"daemon off;\"]"` corresponds to Docker exec form
  - EX: `"CMD nginx -g daemon off;"` corresponds to Docker shell form, invokes a command shell first
- ENTRYPOINT
  - String, supports both array (escaped) and string form
  - EX: `"ENTRYPOINT [\"/bin/sh\", \"-c\", \"/var/www/start.sh\"]"` corresponds to Docker exec form
  - EX: `"ENTRYPOINT /var/www/start.sh"` corresponds to Docker shell form, invokes a command shell first
- ENV
  - String, note there is no equal sign:
  - EX: `"ENV HOSTNAME www.example.com"` not
    `"ENV HOSTNAME=www.example.com"`
- EXPOSE
  - String, space separated ports
  - EX: `"EXPOSE 80 443"`
- LABEL
  - String, space separated key=value pairs
  - EX: `"LABEL version=1.0"`
- ONBUILD
  - String
  - EX: `"ONBUILD RUN date"`
- MAINTAINER
  - String, deprecated in Docker version 1.13.0
  - EX: `"MAINTAINER NAME"`
- USER
  - String
  - EX: `"USER USERNAME"`
- VOLUME
  - String
  - EX: `"VOLUME FROM TO"`
- WORKDIR
  - String
  - EX: `"WORKDIR PATH"`

## Configuration Reference

Configuration options are organized below into two categories: required and
optional. Within each category, the available options are alphabetized and
described.

The Docker builder uses a special Docker communicator _and will not use_ the
standard [communicators](/packer/docs/templates/legacy_json_templates/communicators).

### Required:

You must specify (only) one of `commit`, `discard`, or `export_path`.

<!-- Code generated from the comments of the Config struct in builder/docker/config.go; DO NOT EDIT MANUALLY -->

- `commit` (bool) - If true, the container will be committed to an image rather than exported.
  Default `false`. If `commit` is `false`, then either `discard` must be
  set to `true` or an `export_path` must be provided.

- `discard` (bool) - Throw away the container when the build is complete. This is useful for
  the [artifice
  post-processor](/packer/docs/post-processor/artifice).

- `export_path` (string) - The path where the final container will be exported as a tar file.

- `message` (string) - Set a message for the commit.

<!-- End of code generated from the comments of the Config struct in builder/docker/config.go; -->


### Optional:

<!-- Code generated from the comments of the Config struct in builder/docker/config.go; DO NOT EDIT MANUALLY -->

- `build` (DockerfileBootstrapConfig) - Configuration for a bootstrap image derived from a Dockerfile
  
  Specifying this will make the builder run `docker build` on a provided
  Dockerfile, and this image will then be used to perform the rest of
  the build process.
  
  For more information on the contents of this object, refer to the
  [Bootstrapping a build with a Dockerfile](#bootstrapping-a-build-with-a-dockerfile)
  section of this documentation.

- `author` (string) - Set the author (e-mail) of a commit.

- `changes` ([]string) - Dockerfile instructions to add to the commit. Example of instructions
  are CMD, ENTRYPOINT, ENV, and EXPOSE. Example: [ "USER ubuntu", "WORKDIR
  /app", "EXPOSE 8080" ]

- `container_dir` (string) - The directory inside container to mount temp directory from host server
  for work [file provisioner](/packer/docs/provisioner/file). This defaults
  to c:/packer-files on windows and /packer-files on other systems.

- `device` ([]string) - An array of devices which will be accessible in container when it's run
  without `--privileged` flag.

- `cap_add` ([]string) - An array of additional [Linux
  capabilities](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities)
  to grant to the container.

- `cap_drop` ([]string) - An array of [Linux
  capabilities](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities)
  to drop from the container.

- `docker_path` (string) - Sets the docker binary to use for running commands.
  
  If you want to use a specific version of the docker binary, or a
  docker alternative for building your container, you can specify this
  through this option.
  **Note**: if using an alternative like `podman`, not all options are
  equivalent, and the build may fail in this case.
  
  Defaults to "docker"

- `exec_user` (string) - Username (UID) to run remote commands with. You can also set the group
  name/ID if you want: (UID or UID:GID). You may need this if you get
  permission errors trying to run the shell or other provisioners.

- `image` (string) - The base image for the Docker container that will be started. This image
  will be pulled from the Docker registry if it doesn't already exist.
  Any value format that you can provide to `docker pull` is valid.
  Example: `ubuntu` or `ubuntu:xenial`. If you only provide the repo, Docker
  will pull the latest image, so setting `ubuntu` is the same as setting
  `ubuntu:latest`. You can also set a distribution digest. For example,
  ubuntu@sha256:a0d9e826ab87bd665cfc640598a871b748b4b70a01a4f3d174d4fb02adad07a9
  
  This cannot be used at the same time as `build`

- `privileged` (bool) - If true, run the docker container with the `--privileged` flag. This
  defaults to false if not set.

- `runtime` (string) - Set the container runtime. A runtime different from the one installed
  by default with Docker (`runc`) must be installed and configured.
  The possible values are (non-exhaustive list):
  `runsc` for [gVisor](https://gvisor.dev/),
  `kata-runtime` for [Kata Containers](https://katacontainers.io/),
  `sysbox-runc` for [Nestybox](https://www.nestybox.com/).

- `pull` (bool) - If true, the configured image will be pulled using `docker pull` prior
  to use. Otherwise, it is assumed the image already exists and can be
  used. This defaults to true if not set.
  
  If using `build`, this field will be ignored, as the `pull` option for
  this operation will instead have precedence.

- `run_command` ([]string) - An array of arguments to pass to docker run in order to run the
  container. By default this is set to `["-d", "-i", "-t",
  "--entrypoint=/bin/sh", "--", "{{.Image}}"]` if you are using a linux
  container, and `["-d", "-i", "-t", "--entrypoint=powershell", "--",
  "{{.Image}}"]` if you are running a windows container. `{{.Image}}` is a
  template variable that corresponds to the image template option. Passing
  the entrypoint option this way will make it the default entrypoint of
  the resulting image, so running docker run -it --rm  will start the
  docker image from the /bin/sh shell interpreter; you could run a script
  or another shell by running docker run -it --rm  -c /bin/bash. If your
  docker image embeds a binary intended to be run often, you should
  consider changing the default entrypoint to point to it.

- `tmpfs` ([]string) - An array of additional tmpfs volumes to mount into this container.

- `volumes` (map[string]string) - A mapping of additional volumes to mount into this container. The key of
  the object is the host path, the value is the container path.

- `fix_upload_owner` (bool) - If true, files uploaded to the container will be owned by the user the
  container is running as. If false, the owner will depend on the version
  of docker installed in the system. Defaults to true.

- `windows_container` (bool) - If "true", tells Packer that you are building a Windows container
  running on a windows host. This is necessary for building Windows
  containers, because our normal docker bindings do not work for them.

- `platform` (string) - Set platform if server is multi-platform capable

- `login` (bool) - This is used to login to dockerhub to pull a private base container. For
  pushing to dockerhub, see the docker post-processors

- `login_password` (string) - The password to use to authenticate to login.

- `login_server` (string) - The server address to login to.

- `login_username` (string) - The username to use to authenticate to login.

- `ecr_login` (bool) - Defaults to false. If true, the builder will login in order to pull the
  image from Amazon EC2 Container Registry (ECR). The builder only logs in
  for the duration of the pull. If true login_server is required and
  login, login_username, and login_password will be ignored. For more
  information see the section on ECR.

<!-- End of code generated from the comments of the Config struct in builder/docker/config.go; -->


<!-- Code generated from the comments of the AwsAccessConfig struct in builder/docker/ecr_login.go; DO NOT EDIT MANUALLY -->

- `aws_access_key` (string) - The AWS access key used to communicate with AWS.

- `aws_secret_key` (string) - The AWS secret key used to communicate with AWS.

- `aws_token` (string) - The AWS access token to use. This is different from
  the access key and secret key. If you're not sure what this is, then you
  probably don't need it. This will also be read from the AWS_SESSION_TOKEN
  environmental variable.

- `aws_profile` (string) - The AWS shared credentials profile used to communicate with AWS.

- `aws_force_use_public_ecr` (bool) - The flag to identify whether to push docker image to Public _or_ Private
  ECR. If the user sets this to `true` from the config, we will forcefully
  try to push to Public ECR otherwise set this from code based on the
  given LoginServer value.

<!-- End of code generated from the comments of the AwsAccessConfig struct in builder/docker/ecr_login.go; -->


## Bootstrapping a build with a Dockerfile

The `build` section of a template allows you to specify a Dockerfile to use for bootstrapping a packer build with a locally-built image.

When using this, you won't be able to specify an image as source for the container, instead the image built from the Dockerfile will be used for the remainder of the build after that initial step.

### Configuration examples:

**HCL2**

```hcl
source "docker" "example" {
    build {
        path = "Dockerfile"
    }
    commit = true
}
```

**JSON**

```json
{
  "type": "docker",
  "build": {
    "path": "Dockerfile"
  },
  "commit": true,
}
```

### Required:

<!-- Code generated from the comments of the DockerfileBootstrapConfig struct in builder/docker/dockerfile_config.go; DO NOT EDIT MANUALLY -->

- `path` (string) - Path to the dockerfile to use for building the base image
  
  If set, the builder will invoke `docker build` on it, and use the
  produced image to continue the build afterwards.
  
  Note: Mutually exclusive with "image"

<!-- End of code generated from the comments of the DockerfileBootstrapConfig struct in builder/docker/dockerfile_config.go; -->


### Optional:

<!-- Code generated from the comments of the DockerfileBootstrapConfig struct in builder/docker/dockerfile_config.go; DO NOT EDIT MANUALLY -->

- `build_dir` (string) - Directory to invoke `docker build` from
  
  Defaults to the directory from which we invoke packer.

- `pull` (boolean) - Pull the image when building the base docker image.
  
  Note: defaults to true, to disable this, explicitly set it to false.

- `compress` (bool) - Compress the build context before sending to the docker daemon.
  
  This is especially useful if the build context is large, as copying it
  can take a significant amount of time, while once compressed, this
  can make builds faster, at the price of extra CPU resources.

<!-- End of code generated from the comments of the DockerfileBootstrapConfig struct in builder/docker/dockerfile_config.go; -->


## Build Shared Information Variables

This build shares generated data with provisioners and post-processors via [template engines](/packer/docs/templates/legacy_json_templates/engine)
for JSON and [contextual variables](/packer/docs/templates/hcl_templates/contextual-variables) for HCL2.

The generated variable available for this builder is:

- `ImageSha256` - When committing a container to an image, this will give the image SHA256. Because the image is not available at the provision step,
  this variable is only available for post-processors.

## Using the Artifact: Export

Once the tar artifact has been generated, you will likely want to import, tag,
and push it to a container repository. Packer can do this for you automatically
with the [docker-import](/packer/integrations/hashicorp/docker/latest/components/post-processor/docker-import) and
[docker-push](/packer/integrations/hashicorp/docker/latest/components/post-processor/docker-push) post-processors.

**Note:** This section is covering how to use an artifact that has been
_exported_. More specifically, if you set `export_path` in your configuration.
If you set `commit`, see the next section.

The example below shows a full configuration that would import and push the
created image. This is accomplished using a sequence definition (a collection
of post-processors that are treated as as single pipeline, see
[Post-Processors](/packer/docs/templates/legacy_json_templates/post-processors) for more information):

**JSON**

```json
{
  "post-processors": [
    [
      {
        "type": "docker-import",
        "repository": "myrepo/myimage",
        "tag": "0.7"
      },
      {
        "type": "docker-push"
      }
    ]
  ]
}
```

**HCL2**

```hcl
  post-processors {
    post-processor "docker-import" {
        repository =  "myrepo/myimage"
        tag = "0.7"
      }
    post-processor "docker-push" {}
  }
}
```

In the above example, the result of each builder is passed through the defined
sequence of post-processors starting first with the `docker-import`
post-processor which will import the artifact as a docker image. The resulting
docker image is then passed on to the `docker-push` post-processor which
handles pushing the image to a container repository.

If you want to do this manually, however, perhaps from a script, you can import
the image using the process below:

```shell-session
$ docker import - registry.mydomain.com/mycontainer:latest < artifact.tar
```

You can then add additional tags and push the image as usual with `docker tag`
and `docker push`, respectively.

## Using the Artifact: Committed

If you committed your container to an image, you probably want to tag, save,
push, etc. Packer can do this automatically for you. An example is shown below
which tags and pushes an image. This is accomplished using a sequence
definition (a collection of post-processors that are treated as as single
pipeline, see [Post-Processors](/packer/docs/templates/legacy_json_templates/post-processors) for more
information):

**HCL2**

```hcl
  post-processors {
    post-processor "docker-tag" {
        repository =  "myrepo/myimage"
        tags = ["0.7"]
      }
    post-processor "docker-push" {}
  }
}
```

**JSON**

```json
{
  "post-processors": [
    [
      {
        "type": "docker-tag",
        "repository": "myrepo/myimage",
        "tags": ["0.7"]
      },
      {
        "type": "docker-push"
      }
    ]
  ]
}
```

In the above example, the result of each builder is passed through the defined
sequence of post-processors starting first with the `docker-tag` post-processor
which tags the committed image with the supplied repository and tag
information. Once tagged, the resulting artifact is then passed on to the
`docker-push` post-processor which handles pushing the image to a container
repository.

Going a step further, if you wanted to tag and push an image to multiple
container repositories, this could be accomplished by defining two,
nearly-identical sequence definitions, as demonstrated by the example below:

**HCL2**

```hcl
  post-processors {
    post-processor "docker-tag" {
        repository =  "myrepo/myimage1"
        tags = ["0.7"]
      }
    post-processor "docker-push" {}
  }
  post-processors {
    post-processor "docker-tag" {
        repository =  "myrepo/myimage2"
        tags = ["0.7"]
      }
    post-processor "docker-push" {}
  }
}
```

**JSON**

```json
{
  "post-processors": [
    [
      {
        "type": "docker-tag",
        "repository": "myrepo/myimage1",
        "tags": ["0.7"]
      },
      "docker-push"
    ],
    [
      {
        "type": "docker-tag",
        "repository": "myrepo/myimage2",
        "tags": ["0.7"]
      },
      "docker-push"
    ]
  ]
}
```

<span id="amazon-ec2-container-registry"></span>

## Docker For Windows

You should be able to run docker builds against both linux and Windows
containers. Windows containers use a different communicator than linux
containers, because Windows containers cannot use `docker cp`.

If you are building a Windows container, you must set the template option
`"windows_container": true`. Please note that docker cannot export Windows
containers, so you must either commit or discard them.

The following is a fully functional template for building a Windows
container.

**HCL2**

```hcl
source "docker" "windows" {
    image = "microsoft/windowsservercore:1709"
    container_dir = "c:/app"
    windows_container = true
    commit = true
}

build {
  sources = ["source.docker.example"]
}
```

**JSON**

```json
{
  "builders": [
    {
      "type": "docker",
      "image": "microsoft/windowsservercore:1709",
      "container_dir": "c:/app",
      "windows_container": true,
      "commit": true
    }
  ]
}
```

## Amazon EC2 Container Registry

Packer can tag and push images for use in [Amazon EC2 Container
Registry](https://aws.amazon.com/ecr/). The post processors work as described
above and example configuration properties are shown below:

**HCL2**

```hcl
post-processors {
  post-processor "docker-tag" {
    repository =  "12345.dkr.ecr.us-east-1.amazonaws.com/packer"
    tags       = ["0.7"]
  }

  post-processor "docker-push" {
    ecr_login = true
    aws_access_key = "YOUR KEY HERE"
    aws_secret_key = "YOUR SECRET KEY HERE"
    login_server = "https://12345.dkr.ecr.us-east-1.amazonaws.com/"
  }
}
```

**JSON**

```json
{
  "post-processors": [
    [
      {
        "type": "docker-tag",
        "repository": "12345.dkr.ecr.us-east-1.amazonaws.com/packer",
        "tags": ["0.7"]
      },
      {
        "type": "docker-push",
        "ecr_login": true,
        "aws_access_key": "YOUR KEY HERE",
        "aws_secret_key": "YOUR SECRET KEY HERE",
        "login_server": "https://12345.dkr.ecr.us-east-1.amazonaws.com/"
      }
    ]
  ]
}
```

## Amazon ECR Public Gallery

Packer can tag and push images for use in [Amazon ECR Public
Gallery](https://gallery.ecr.aws/). The post processors work as described above
and example configuration properties are shown below:

**HCL2**

```hcl
    post-processors {
        post-processor "docker-tag" {
            repository = "public.ecr.aws/YOUR REGISTRY ALIAS HERE/YOUR REGISTRY NAME HERE"
            tags       = ["latest"]
        }

        post-processor "docker-push" {
            "ecr_login": true,
            "aws_access_key": "YOUR KEY HERE",
            "aws_secret_key": "YOUR SECRET KEY HERE",
            login_server = "public.ecr.aws/YOUR REGISTRY ALIAS HERE"
        }
    }
```

**JSON**

```json
    {
        "post-processors": [
            [
                {
                    "type": "docker-tag",
                    "repository": "public.ecr.aws/YOUR REGISTRY ALIAS HERE/YOUR REGISTRY NAME HERE",
                    "tags": ["latest"]
                },
                {
                    "type": "docker-push",
                    "ecr_login": true,
                    "aws_access_key": "YOUR KEY HERE",
                    "aws_secret_key": "YOUR SECRET KEY HERE",
                    "login_server": "public.ecr.aws/YOUR REGISTRY ALIAS HERE"
                }
            ]
        ]
    }
```

[Learn how to set Amazon AWS credentials.](/packer/integrations/hashicorp/amazon#specifying-amazon-credentials)

## Dockerfiles

This builder allows you to build Docker images _without_ Dockerfiles.

With this builder, you can repeatedly create Docker images without the use of a
Dockerfile. You don't need to know the syntax or semantics of Dockerfiles.
Instead, you can just provide shell scripts, Chef recipes, Puppet manifests,
etc. to provision your Docker container just like you would a regular
virtualized or dedicated machine.

While Docker has many features, Packer views Docker simply as a container
runner. To that end, Packer is able to repeatedly build these containers using
portable provisioning scripts.

**Note**: starting with v1.1.0 of the plugin, this builder supports bootstrapping
a build from a Dockerfile. This slightly conflicts with the original intent, but
practically, for users who already have working Dockerfile-centric pipelines,
this limitation was a hinderance to adopting Packer for later provisioning images,
so we opted to add this capability to the builder.

## Overriding the host directory

By default, Packer creates a temporary folder under your home directory, and
uses that to stage files for uploading into the container. If you would like to
change the path to this temporary folder, you can set the `PACKER_TMP_DIR`.
This can be useful, for example, if you have your home directory permissions
set up to disallow access from the docker daemon.
