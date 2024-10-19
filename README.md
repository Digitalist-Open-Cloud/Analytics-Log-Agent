# Matomo Agent

This is an agent for Matomo, running on the server you want to have Matomo data from, to send to your Matomo instance, like 404 logs (this needs the Matomo plugin `Agent` installed and activated in you Matomo instance).

The agent is in it's early stages, and at this point **NOT** recommended for production use.

The idea of the usages of this agent are:

- If you can't use normal JavaScript tracking for you website, you instead use the logs and tail them and send the data to Matomo. (If you want to send server logs by daily batches, you instead use [Log Analytics](https://matomo.org/log-analytics/)).
- Track user agents that doesn't use a browser to visit your sites - like some bots and people using command lines tool to scrap your website.
- Generate test data for you Matomo with logs generated from something like [Dummy Log Creator](https://github.com/Digitalist-Open-Cloud/dummy-log-generator).

For now, you need to build the agent manually, and then you can add in your preferred path. A config file is required, and as default it should be placed `/opt/matomo-agent/config.toml`, you can override this when starting the Matomo agent with the `-config` flag pointing to path to the file, like:

```sh
./matomo-agent --config /usr/local/bin/matomo-agent-config.toml
```

## Config

Options for `config.toml`:

| Config              | Description                                   | Default | Required |
| ------------------- | --------------------------------------------- | ------- | -------- |
| `matomo.url`        | URL to your Matomo instances                  | -       | Yes      |
| `matomo.site_id`    | Site id in Matomo to track to                 | 1       | Yes      |
| `matomo.token_auth` | Token auth to your Matomo instance            | -       | Yes      |
| `matomo.plugin`     | If you want to use the Agent plugin in Matomo | false   | No       |
| `log.log_format`    | Which log format the log has                  | -       | Yes      |
| `log.log_path`      | Path to the log to tail                       | -       | Yes      |
| `agent.log_level`   | Log level for Matomo agent                    | -       | Yes      |
| `agent.log_file`    | File to log to                                | -       | Yes      |

## Log format

The plan is to support different logging formats, that is why apache and nginx
are added as formats, but for now the config is the same, as the apache and
nginx format supported now is the common combined log format.

## Build

go build -o matomo-agent .

## Install

Copy config.toml.example to default, /opt/matomo-agent/config.toml
Add settings for the agent in config.toml, run matomo-agent.

```sh
./matomo-agent
```

`config.toml` could also be placed in another place, then you need to start
the agent with the `--config` flag.

```sh
./matomo-agent --config config.toml
```

## Todos

- If activated config, send information about 404, 503 etc. to Matomo - using the Matomo plugin Agent as API endpoint. Only send ok responses to Matomo Tracking API.
- If download, add metadata to Tracking API.
- To configure to track bots only, this could be used: <https://github.com/robicode/device-detector/tree/main>
- Add releases <https://github.com/goreleaser/goreleaser>
