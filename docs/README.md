# Docker Plugins

The Docker plugin is intended to be used for managing docker images through Packer:

## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    docker = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/docker"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/hashicorp/packer-plugin-docker/releases).
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


#### From Source

If you prefer to build the plugin from its source code, clone the GitHub
repository locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-docker` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).

## Builders:
- [docker](/docs/builders/docker.mdx) - The builder builds Docker images using Docker.
  The builder starts a Docker container, runs provisioners within this container, then exports the container for reuse or commits the image.

## Post-Processors
- [docker-import](/docs/post-processors/docker-import.mdx) - The import post-processor takes an artifact from the docker builder and imports it with Docker locally.
- [docker-push](/docs/post-processors/docker-push.mdx) - The push post-processor takes an artifact from the docker-import post-processor and pushes it to a Docker registry.
- [docker-save](/docs/post-processors/docker-save.mdx) - The save post-processor takes an artifact from the docker builder that was committed and saves it to a file.
- [docker-tag](/docs/post-processors/docker-tag.mdx) - The tag post-processor takes an artifact from the docker builder that was committed and tags it into a repository.

