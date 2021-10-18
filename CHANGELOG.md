## 1.0.2 (October 18, 2021)

### IMPROVEMENTS:
* Add `SourceImageDigest` and `ImageSha256` as shared builder information.
    [GH-75].
* Update plugin to Go 1.17
* Update packer-plugin-sdk to v0.2.7 [GH-]
* Small refactor to main driver to support the capturing of the image digest
    for the source image. [GH-75]

## 1.0.1 (June 15, 2021)

`1.0.1` is the same as `1.0.0`

## 1.0.0 (June 14, 2021)
* Update packer-plugin-sdk to v0.2.3. [GH-56]

## 0.0.7 (March 30, 2021)
* Docker plugin break out from Packer core. Changes prior to break out can be found in [Packer's CHANGELOG](https://github.com/hashicorp/packer/blob/master/CHANGELOG.md)

### BUG FIXES
* Update packer-plugin-sdk to latest version to fix issue with file provisioner. [GH-24]
