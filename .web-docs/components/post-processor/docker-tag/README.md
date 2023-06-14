Type: `docker-tag`

The Packer Docker Tag post-processor takes an artifact from the [docker
builder](/packer/integrations/hashicorp/docker) that was committed and tags it into a
repository. This allows you to use the other Docker post-processors such as
[docker-push](/packer/integrations/hashicorp/docker/latest/components/post-processor/docker-push) to push the image to a
registry.

This is very similar to the
[docker-import](/packer/integrations/hashicorp/docker/latest/components/post-processor/docker-import) post-processor except
that this works with committed resources, rather than exported.

## Configuration

The configuration for this post-processor requires `repository`, all other
settings are optional.

- `repository` (string) - The repository of the image.

- `tags` (array of strings) - A list of tags for the image. By default this is
  not set. Example of declaration: `"tags": ["mytag-1", "mytag-2"]`

- `force` (boolean) - If true, this post-processor forcibly tag the image
  even if tag name is collided. Default to `false`. But it will be ignored if
  Docker &gt;= 1.12.0 was detected, since the `force` option was removed
  after 1.12.0.
  [reference](https://docs.docker.com/engine/deprecated/#/f-flag-on-docker-tag)

- `keep_input_artifact` (boolean) - Unlike most other post-processors, the
  keep_input_artifact option will have no effect for the docker-tag
  post-processor. We will always retain the input artifact for docker-tag,
  since deleting the image we just tagged is not a behavior anyone should ever
  expect. `keep_input_artifact will` therefore always be evaluated as true,
  regardless of the value you enter into this field.

## Example

An example is shown below, showing only the post-processor configuration:

**JSON**

```json
{
  "type": "docker-tag",
  "repository": "hashicorp/packer",
  "tags": ["0.7", "anothertag"]
}
```

**HCL2**

```hcl
post-processor "docker-tag" {
  repository = "hashicorp/packer"
  tags = ["0.7", "anothertag"]
}
```

This example would take the image created by the Docker builder and tag it into
the local Docker process with a name of `hashicorp/packer:0.7`.

Following this, you can use the
[docker-push](/packer/integrations/hashicorp/docker/latest/components/post-processor/docker-push) post-processor to push it
to a registry, if you want.
