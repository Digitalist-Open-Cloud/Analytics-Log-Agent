# Matomo Agent

This is an agent for Matomo, running on the server you want to have Matomo data from, to send to your Matomo instance.

There is also an integration for the Matomo plugin `Agent` that could be used for cases when you want to have reports for errors, like 404 logs.

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

The file `config.toml.example` kan be used as a start for your `config.tml` file.

## Config

### Flags

| **Flag**        | **Type** | **Default**                     | **Description**                                                                                   |
| --------------- | -------- | ------------------------------- | ------------------------------------------------------------------------------------------------- |
| `--config`      | `string` | `/opt/matomo-agent/config.toml` | Path to the configuration file.                                                                   |
| `--catlog`      | `bool`   | `false`                         | Simulate `cat` command for a log file. If set to `true`, processes log file in one go.            |
| `--rps`         | `int`    | `1`                             | Requests per second limit for `catlog` mode. Controls the rate of log file processing.            |
| `--matomo-url`  | `string` | `""`                            | Matomo URL. Overrides the value set in the config file.                                           |
| `--token-auth`  | `string` | `""`                            | Matomo authentication token. Overrides the value set in the config file.                          |
| `--site-id`     | `string` | `""`                            | Matomo site ID. Overrides the value set in the config file.                                       |
| `--plugin`      | `bool`   | `false`                         | If using the Matomo Agent plugin, set this flag to enable plugin functionality.                   |
| `--downloads`   | `bool`   | `true`                          | Enable or disable download tracking. Overrides the config file setting.                           |
| `--log-format`  | `string` | `""`                            | Log format. Valid options: `nginx`, `apache`, or `csv`. Overrides the config file setting.        |
| `--log-path`    | `string` | `""`                            | Path to the log file. Overrides the value set in the config file.                                 |
| `--user-agents` | `string` | `""`                            | Comma-separated list of user agents to track. Overrides the config file setting.                  |
| `--log-level`   | `string` | `""`                            | Log level. Valid options: `debug`, `info`, `warn`, or `error`. Overrides the config file setting. |
| `--log-file`    | `string` | `""`                            | Path to the agent's log file. Overrides the value set in the config file.                         |

Each flag can be used to override corresponding values in the `config.toml` file, allowing you to customize the agent's behavior via command-line arguments.

### File

Options for `config.toml`:

| Config              | Description                                   | Default | Required |
| ------------------- | --------------------------------------------- | ------- | -------- |
| `matomo.url`        | URL to your Matomo instances                  | -       | Yes      |
| `matomo.site_id`    | Site id in Matomo to track to                 | 1       | Yes      |
| `matomo.token_auth` | Token auth to your Matomo instance            | -       | Yes      |
| `matomo.plugin`     | If you want to use the Agent plugin in Matomo | false   | No       |
| `matomo.downloads`  | If you want to track downloads                | true    | No       |
| `log.log_format`    | Which log format the log has                  | -       | Yes      |
| `log.log_path`      | Path to the log to tail                       | -       | Yes      |
| `log.user_agents`   | Array of User Agents that should be tracked   | -       | No       |
| `agent.log_level`   | Log level for Matomo agent                    | -       | Yes      |
| `agent.log_file`    | File to log to                                | -       | Yes      |

## Log format

### Apache and Nginx

Apache and Nginx log format supported is the combined log format, that matches this regex:

```sh
(?P<ip>\S+) - - \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<url>\S+) (?P<protocol>\S+)" (?P<status>\d+) (?P<size>\d+) "(?P<referrer>[^"]*)" "(?P<user_agent>[^"]*)"
```

### CSV

If using CSV file, the format need to be:

```sh
timestamp, req_method, req_host, req_uri, resp_status, client_ip, req_referer, req_user_agent
```

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

## Tail or cat

### Tail

Just start like:

```sh
./matomo-agent [--config config.toml]
```

### Cat

As an option you could read the logfile from start to end, if you have a log file that is not updated anymore, to do that you could run it the agent like:

```sh
./matomo-agent --config config.toml --catlog --rps 5
```

The `catlog` flag makes the agent run just once, and the `rps` flag is to set how many requests per second if needed, default is `1`.

We do though recommend using Matomos official Log Analytics for this.

## Todos

- To configure to track bots only, this could be used: <https://github.com/robicode/device-detector/tree/main>
- Add releases <https://github.com/goreleaser/goreleaser>

## Linting

```sh
golangci-lint run
```

## License

Copyright (C) 2024 Digitalist Open Cloud <cloud@digitalist.com>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>
