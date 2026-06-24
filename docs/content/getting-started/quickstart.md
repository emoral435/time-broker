---
title: Quick Start
weight: 2
---

## 1. Install

```shell
curl -fsSL https://raw.githubusercontent.com/emoral435/time-broker/main/install.sh | sh
```

This installs `time-broker` and the `tb` shorthand to your PATH.

## 2. Set up

Upon first run of your broker, you will be presented with a wizard to configure your brokers settings.
```shell
time-broker
```

Follow the interactive wizard to choose your calendar provider and week start day.

## 3. Check your config

The following will show you what settings are possible to configure.
```shell
time-broker config list
```

## 4. View available commands

```shell
time-broker help
```
