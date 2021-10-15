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

## Ci System Support

We build a [Container Image](https://github.com/docat-org/docat-cli/pkgs/container/docatl) you can use
in your Ci system.

### GitLab Ci

Use the following Job template to publish the docs:

```sh
deploy-docs:
  image: ghcr.io/docat-org/docatl:latest
  script:
    - --host docat.company.io ./docs.zip $CI_PROJECT $CI_COMMIT_TAG
```