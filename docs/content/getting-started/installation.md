---
title: Installation
weight: 1
---

## Using curl

```shell
curl -fsSL https://raw.githubusercontent.com/emoral435/time-broker/main/install.sh | sh
```

This downloads the latest release for your OS and architecture, installs it to
`~/.time-broker/bin/`, and creates a `tb` shorthand symlink.

After installing, use `time-broker` or `tb` from your terminal.

## From a Release

Download the latest binary from the
[releases page](https://github.com/emoral435/time-broker/releases) for your
platform:

- macOS (arm64 / amd64)
- Linux (arm64 / amd64)
- Windows (arm64 / amd64)

## Uninstall

To remove time-broker, run:

```sh
time-broker uninstall
```

This will prompt you for confirmation and then remove the installed binaries,
configuration directory (`~/.time-broker/`), and any symlinks that were created
during installation.
