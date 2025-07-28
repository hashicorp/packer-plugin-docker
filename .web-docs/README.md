The Docker plugin is intended to be used for managing docker images through Packer.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    docker = {
      source  = "github.com/hashicorp/docker"
      version = "~> 1"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/hashicorp/docker
```

**Note: Update to Packer Plugin Installation**

With the new Packer release starting from version 1.14.0, the `packer init` command will automatically install official plugins from the [HashiCorp release site.](https://releases.hashicorp.com/)

Going forward, to use newer versions of official Packer plugins, you'll need to upgrade to Packer version 1.14.0 or later. If you're using an older version, you can still install plugins, but as a workaround, you'll need to [manually install them using the CLI.](https://developer.hashicorp.com/packer/docs/plugins/install#manually-install-plugins-using-the-cli)

There is no change to the syntax or commands for installing plugins.

### Components

#### Builders

- [docker](/packer/integrations/hashicorp/docker/latest/components/builder/docker) - The builder builds Docker images using Docker.
  The builder starts a Docker container, runs provisioners within this container, then exports the container for reuse or commits the image.

#### Post-Processors

- [docker-import](/packer/integrations/hashicorp/docker/latest/components/post-processor/docker-import) - The import post-processor
  takes an artifact from the docker builder and imports it with Docker locally.

- [docker-push](/packer/integrations/hashicorp/docker/latest/components/post-processor/docker-push) - The push post-processor takes
  an artifact from the docker-import post-processor and pushes it to a Docker registry.

- [docker-save](/packer/integrations/hashicorp/docker/latest/components/post-processor/docker-save) - The save post-processor takes
  an artifact from the docker builder that was committed and saves it to a file.

- [docker-tag](/packer/integrations/hashicorp/docker/latest/components/post-processor/docker-tag) - The tag post-processor takes an
  artifact from the docker builder that was committed and tags it into a repository.
