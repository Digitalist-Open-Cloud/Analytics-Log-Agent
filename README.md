# Matomo Agent

This is an agent for Matomo, running on the server you want to have Matomo data from, to send to your Matomo instance, like 404 logs (this needs the Matomo plugin `Agent` installed and activated in you Matomo instance).

The agent is in it's early stages, and at this point **NOT** recommended for production use.

## Usage

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
| `log.user_agents`   | Array of User Agents that should be tracked   | -       | No      |
| `agent.log_level`   | Log level for Matomo agent                    | -       | Yes      |
| `agent.log_file`    | File to log to                                | -       | Yes      |

## Log format

The plan is to support different logging formats, that is why apache and nginx
are added as formats, but for now the config is the same, as the apache and
nginx format supported now is the common combined log format.

## Build

go build -o matomo-agent .

## Install

Copy config.toml.example to default (are you preferred destination), /opt/matomo-agent/config.toml
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

## Linting

```sh
golangci-lint run
```

## License

Copyright (C) 2024 Digitalist Open Cloud <cloud@digitalist.com>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program.  If not, see <https://www.gnu.org/licenses/>