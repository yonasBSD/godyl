---
layout: default
title: Commands
---

# Command Reference

`godyl` provides several commands to help you manage your tools. This section provides detailed information about each command and its options.

## Available Commands

| Command                                                          | Description                           |
| ---------------------------------------------------------------- | ------------------------------------- |
| [`install`]({{ site.baseurl }}/commands/install#content-start)   | Install tools from YAML files         |
| [`download`]({{ site.baseurl }}/commands/download#content-start) | Download and install individual tools |
| [`status`]({{ site.baseurl }}/commands/status#content-start)     | Check the status of installed tools   |
| [`dump`]({{ site.baseurl }}/commands/dump#content-start)         | Display configuration information     |
| [`update`]({{ site.baseurl }}/commands/update#content-start)     | Update the godyl application          |
| [`cache`]({{ site.baseurl }}/commands/cache#content-start)       | Manage the godyl cache                |
| [`version`]({{ site.baseurl }}/commands/version#content-start)   | Display the current version of godyl  |

## Global Flags

These flags are available for all commands:

| Flag                  | Environment Variable     | Default                              | Description                                          |
| --------------------- | ------------------------ | ------------------------------------ | ---------------------------------------------------- |
| `--help`, `-h`        | `GODYL_HELP`             | `false`                              | Show help message and exit                           |
| `--version`           | `GODYL_VERSION`          | `false`                              | Show version information and exit                    |
| `--dry`               | `GODYL_DRY`              | `false`                              | Run without making any changes (dry run)             |
| `--log`               | `GODYL_LOG`              | `info`                               | Log level (debug, info, warn, error, silent)         |
| `--config-file`, `-c` | `GODYL_CONFIG_FILE`      | `${XDG_CONFIG_HOME}/godyl/godyl.yml` | Path to `godyl.yml` file                             |
| `--env-file`, `-e`    | `GODYL_ENV_FILE`         | `[".env"]`                           | Path to `.env` file(s)                               |
| `--defaults`, `-d`    | `GODYL_DEFAULTS`         | `defaults.yml`                       | Path to defaults file                                |
| `--github-token`      | `GODYL_GITHUB_TOKEN`     | `${GODYL_GITHUB_TOKEN}`              | GitHub token for authentication                      |
| `--gitlab-token`      | `GODYL_GITLAB_TOKEN`     | `${GODYL_GITLAB_TOKEN}`              | GitLab token for authentication                      |
| `--url-token`         | `GODYL_URL_TOKEN`        | `${GODYL_URL_TOKEN}`                 | URL token for authentication                         |
| `--url-token-header`  | `GODYL_URL_TOKEN_HEADER` | `Authorization`                      | URL header for authentication                        |
| `--cache-dir`         | `GODYL_CACHE_DIR`        | `${XDG_CACHE_HOME}/godyl`            | Path to cache directory                              |
| `--no-cache`          | `GODYL_NO_CACHE`         | `false`                              | Disable cache usage                                  |
| `--default`           | `GODYL_DEFAULT`          | `default`                            | The default configuration to use from `defaults.yml` |

For `--env-file` and `--defaults`, the defaults are used only if no issue is encountered while loading them.

In addition, the following environment variables will be read directly from the environment (i.e not the `.env` file(s)):

- `--github-token` from `GITHUB_TOKEN` or `GH_TOKEN`
- `--gitlab-token` from `GITLAB_TOKEN`
- `--url-token` from `URL_TOKEN`

The value for `--config-file` defaults to `$XDG_CONFIG_HOME/godyl/godyl.yml`. If `XDG_CONFIG_HOME` is empty, it will be set to `$HOME/.config/godyl/godyl.yml`.
If neither are available, it will default to `./godyl.yml`.

The value for `--cache-dir` defaults to `$XDG_CACHE_HOME/godyl`. If `XDG_CACHE_HOME` is empty, it will be set to `$HOME/.cache/godyl`.
If neither are available, it will default to the subdirectory `godyl` in the system temporary directory.

Temporary assets are downloaded by default in `$XDG_RUNTIME_DIR/godyl`. If `XDG_RUNTIME_DIR` is empty, it will be set to `/tmp/${USER}/godyl`.
If user is not set, it will default to the subdirectory `godyl` in the system temporary directory.

Equivalent fallbacks are made for other platforms such as `Windows` and `macOS`.
