<p align="center">
  <img alt="golangci-lint logo" src="docs/assets/images/godyl.png" height="150" />
  <h3 align="center">godyl</h3>
  <p align="center">Asset downloader</p>
</p>

---

[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/godyl.svg)](https://pkg.go.dev/github.com/idelchi/godyl)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/godyl)](https://goreportcard.com/report/github.com/idelchi/godyl)
[![Build Status](https://github.com/idelchi/godyl/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/godyl/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`godyl` helps with batch-downloading and installing statically compiled binaries from:

- GitHub releases
- GitLab releases
- URLs
- Go projects

As an alternative to above, custom commands can be used as well.

`godyl` will infer the platform and architecture from the system it is running on, and will attempt to download the appropriate binary.

This uses simple heuristics to select the correct binary to download, and will not work for all projects.

However, most properties can be overridden, with `hints` and `skip` used to help `godyl` make the correct decision.

> [!NOTE]
> Set up a GitHub API token to avoid rate limiting when using `github` as a source type.
> See [configuration](#configuration) for more information, or simply `export GODYL_GITHUB_TOKEN=<token>`

Tool is inspired by [task](https://github.com/go-task/task), [dra](https://github.com/devmatteini/dra) and [ansible](https://github.com/ansible/ansible)

For full documentation, see [godyl documentation](https://idelchi.github.io/godyl/)
