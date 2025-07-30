## 1.1.2 (July 31, 2025)

### IMPROVEMENTS:

* core: Changes to pull official packer plugins binaries from official site (releases.hashicorp.com).
  This change allows Packer to automatically download and install official plugins from the HashiCorp official release site.
  This change standardizes our release process and ensures a more secure and reliable pipeline for plugin delivery.
  [GH-13431](https://github.com/hashicorp/packer/pull/13431)

* CRT migration changes by @anshulsharma-hashicorp in https://github.com/hashicorp/packer-plugin-docker/pull/210
  Change in the release process of the packer plugins binaries to releases it in the [HashiCorp official releases site](https://releases.hashicorp.com/packer-plugin-docker/).
  This change standardizes our release process and ensures a more secure and reliable pipeline for plugin delivery.

### Other Changes
* build(deps): bump github.com/hashicorp/packer-plugin-sdk from 0.5.4 to 0.6.0 by @dependabot[bot] in https://github.com/hashicorp/packer-plugin-docker/pull/200
* build(deps): bump github.com/hashicorp/packer-plugin-sdk from 0.6.0 to 0.6.1 by @dependabot[bot] in https://github.com/hashicorp/packer-plugin-docker/pull/203
* build(deps): bump github.com/hashicorp/packer-plugin-sdk from 0.6.1 to 0.6.2 by @dependabot[bot] in https://github.com/hashicorp/packer-plugin-docker/pull/207
* Update PR template for PCI by @devashish-patel in https://github.com/hashicorp/packer-plugin-docker/pull/206
* Log docker build from Dockerfile for debugging by @radtriste in https://github.com/hashicorp/packer-plugin-docker/pull/204
* [COMPLIANCE] Add Copyright and License Headers by @hashicorp-copywrite[bot] in https://github.com/hashicorp/packer-plugin-docker/pull/211
* Manifest json change by @anshulsharma-hashicorp in https://github.com/hashicorp/packer-plugin-docker/pull/212
* typo fix by @anshulsharma-hashicorp in https://github.com/hashicorp/packer-plugin-docker/pull/213

## New Contributors
* @radtriste made their first contribution in https://github.com/hashicorp/packer-plugin-docker/pull/204
* @anshulsharma-hashicorp made their first contribution in https://github.com/hashicorp/packer-plugin-docker/pull/210
---

# Changelog of previous releases can be found below.

Please refer to [releases](https://github.com/hashicorp/packer-plugin-docker/releases) for the latest CHANGELOG information.

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
