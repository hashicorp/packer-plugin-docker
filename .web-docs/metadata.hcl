# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Docker"
  description = "The docker plugin can be used with HashiCorp Packer to manage containers with Docker."
  identifier = "packer/hashicorp/docker"
  component {
    type = "builder"
    name = "Docker"
    slug = "docker"
  }
  component {
    type = "post-processor"
    name = "Docker Import"
    slug = "docker-import"
  }
  component {
    type = "post-processor"
    name = "Docker Save"
    slug = "docker-save"
  }
  component {
    type = "post-processor"
    name = "Docker Tag"
    slug = "docker-tag"
  }
  component {
    type = "post-processor"
    name = "Docker Push"
    slug = "docker-push"
  }
}
