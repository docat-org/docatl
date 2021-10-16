# docatl, the docat cli

![docat](https://github.com/randombenj/docat/raw/master/doc/assets/docat-teaser.png)

**Manage your docat documentation with ease.**

[![build](https://github.com/docat-org/docat-cli/workflows/Ci/badge.svg)](https://github.com/docat-org/docat-cli/actions)

## Getting Started

Download the [latest Release binary](https://github.com/docat-org/docat-cli/releases/latest) for your platform
and start pushing your documentation:

```sh
docatl push --host docat.company.io ./docs.zip myproject v1.0.0
```

## Installation

* Binaries for your platform are attached to each release [here](https://github.com/docat-org/docat-cli/releases)
* Container images are available on ghcr.io [here](https://github.com/docat-org/docat-cli/pkgs/container/docatl)

### Using Go

You can install the package using go directly:

```sh
go install github.com/docat-org/docat-cli
```

The `docatl` binary will be placed in `$GOPATH/bin`.

### Using wget / curl

You can use `wget` or `curl` to download a release.
Make sure that you are using the correct version and platform names in the URL.

E.g downloading `v0.1.0` for `Linux` `x86_64`:
```sh
wget https://github.com/docat-org/docat-cli/releases/download/v0.1.0/docatl_0.1.0_Linux_x86_64 -O ~/bin/docatl
```

## Ci System Support

We build a [Container Image](https://github.com/docat-org/docat-cli/pkgs/container/docatl) you can use
in your Ci system.

### GitLab Ci

Use the following Job template to publish the docs:

```yaml
deploy-docs:
  image: ghcr.io/docat-org/docatl:latest
  variables:
    DOCATL_HOST: https://docat.company.io
    DOCATL_API_KEY: blabla
  script:
    - push ./docs.zip $CI_PROJECT $CI_COMMIT_TAG
```

## Configuration

You can configure `docatl` with either command line arguments, env variables and/or a config file (evaluated in that order of precedence).

### Config File

The config file must be at the working directory at `.docatl.yaml`, e.g.:

```yaml
host: https://docat.company.io
api-key: blabla
```

### Environment Variable

The `DOCATL_` must be used, e.g.:

```sh
DOCATL_API_KEY=blabla docatl push ...
```