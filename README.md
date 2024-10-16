# Readme

## Log format

The plan is to support different logging formats, that is why apache and nginx
are added as formats, but for now the config is the same, as the apache and
nginx format supported is the common combined log format.

## Build

go build -o matomo-agent .

## Install

Copy config.toml.example to default, /opt/matomo-agent/config.toml
Add settings for the agent in config.toml, run matomo-agent.

```sh
./matomo-agent
```

`config.toml` could also be placed in another place, then you need to start
the agent with the `-config` flag.

```sh
./matomo-agent -config /Users/mikke/dev/lab/matomo-agent/config.toml
```
