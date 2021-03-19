# Docker Plugins

The Docker plugin is intended to be used for managing docker images through Packer:

### Builders:
- [docker](/docs/builders/docker.mdx) - The builder builds Docker images using Docker.
  The builder starts a Docker container, runs provisioners within this container, then exports the container for reuse or commits the image.

### Post-Processors
- [docker-import](/docs/post-processors/docker-import.mdx) - The import post-processor takes an artifact from the docker builder and imports it with Docker locally.
- [docker-push](/docs/post-processors/docker-push.mdx) - The push post-processor takes an artifact from the docker-import post-processor and pushes it to a Docker registry.
- [docker-save](/docs/post-processors/docker-save.mdx) - The save post-processor takes an artifact from the docker builder that was committed and saves it to a file.
- [docker-tag](/docs/post-processors/docker-tag.mdx) - The tag post-processor takes an artifact from the docker builder that was committed and tags it into a repository.

