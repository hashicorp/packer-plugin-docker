Type: `docker-save`

The Packer Docker Save post-processor takes an artifact from the [docker
builder](/packer/integrations/hashicorp/docker) that was committed and saves it to a file.
This is similar to exporting the Docker image directly from the builder, except
that it preserves the hierarchy of images and metadata.

We understand the terminology can be a bit confusing, but we've adopted the
terminology from Docker, so if you're familiar with that, then you'll be
familiar with this and vice versa.

## Configuration

### Required

The configuration for this post-processor only requires one option.

- `path` (string) - The path to save the image.

### Optional

- `keep_input_artifact` (boolean) - if true, do not delete the docker
  container, and only save the .tar created by docker save. Defaults to true.

## Example

An example is shown below, showing only the post-processor configuration:

**JSON**

```json
{
  "type": "docker-save",
  "path": "foo.tar"
}
```

**HCL2**

```hcl
post-processor "docker-save" {
  path = "foo.tar"
}
```
