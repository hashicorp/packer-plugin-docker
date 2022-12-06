# Latest Release

Please refer to [releases](https://github.com/hashicorp/packer-plugin-vsphere/releases) for the latest CHANGELOG information.

---
## 1.0.3 (October 29, 2021)

### IMPROVEMENTS:
* Add ImageDigest label to HCP Packer registry image metadata, if available. [GH-79] [GH-80]

### BUG FIXES:
* Properly set ImageID and related HCP Packer registry labels when set. [GH-78]

## 1.0.2 (October 18, 2021)

### NOTES:
Support for the HCP Packer registry is currently in beta and requires
Packer v1.7.7 [GH-75]

### FEATURES:
* Add HCP Packer registry image metadata for all artifacts. [GH-75]

### IMPROVEMENTS:
* Add `SourceImageDigest` and `ImageSha256` as shared builder information.
    [GH-75].
* Update plugin to Go 1.17
* Update packer-plugin-sdk to v0.2.7 [GH-74]
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
