# Emby Exporter

An exporter for emby that exports Emby's main metrics.

## /!\ IMPORTANT /!\

This project may no longer be maintained. I no longer use emby because of the 25-device limit, which I have unfortunately exceeded. 

As I couldn't find a solution with the Emby team and couldn't afford to take out a subscription just to extend the limit, I was forced to stop using emby.

If you ever want to add features, I invite you to create a fork of the project and either continue it on your own, or make a PR that I'll validate if I can.

I will continue to respond to issues and try to solve your problems or improve it as much as possible.

Thank you for your understanding.

### Endpoints

| Endpoint | Description |
|---|---|
| `http://ip:9210/metrics` | Prometheus metrics |
| `http://ip:9210/health` | Returns `200` while the exporter is running (can be disabled with `options.healthcheck: false`) |

### Metrics

| Metric | Description |
|---|---|
| `emby_up` | `1` if the Emby server is reachable, `0` otherwise |
| `emby_system_info` | Emby version, WAN/LAN addresses, pending update/restart |
| `emby_sessions` | Active playback sessions (user, client, media, play method, location…) |
| `emby_sessions_count` | Number of active playback sessions |
| `emby_sessions_bitrate` | Bitrate of each session, in kilobits |
| `emby_media_item` | Number of items in each library |
| `emby_activity` | Latest entries of the activity log |
| `emby_alert` | Latest system alerts |

The exporter also exposes the standard `go_*` and `process_*` metrics (memory, goroutines, CPU…).

When the Emby server is unreachable, `/metrics` only returns `emby_up 0` and the exporter retries every `retryInterval` seconds until the server answers again. You can alert on it with:

```yaml
- alert: EmbyDown
  expr: emby_up == 0
  for: 2m
```

## Grafana Dashboard

![Dashboard example](https://github.com/TOomaAh/emby_exporter_go/blob/main/example/dashboard_grafana.png)

[Dashboard link](https://github.com/TOomaAh/emby_exporter_go/blob/main/example/Emby.Dashboard-1703419734858.json)

I thank [jaycedk](https://github.com/jaycedk) for the dashboard (it's his)


## Build this project

```
go build -o emby_exporter ./cmd/emby_exporter_go
```

## Command line options and environment variables

| Flag | Environment variable | Default | Description |
|---|---|---|---|
| `-c`, `--config` | `CONFIG_FILE` | `./config.yml` | Path of the configuration file |
| `-l`, `--log-level` | `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `-g`, `--geoip` | `GEOIP_DB` | `geoip.mmdb` | Path of the GeoIP database file |
| `-v`, `--version` | | | Print the version and exit |
| | `TZ` | `UTC` | Time zone used for the dates of activities and alerts (ex: `Europe/Paris`) |

A command line flag takes precedence over its environment variable.

## Config file

Only `server.token` is required, every other key has a default value.

```yaml
server:
  url: "http://<ip|domain name>" # default: http://localhost
  port: 8096                     # default: 8096
  token: "your token"            # required
exporter:
  port: 9210                     # default: 9210
options:
  retryInterval: 10              # seconds between two pings while Emby is unreachable, default: 10
  healthcheck: true              # expose /health, default: true
  geoip: false                   # enable GeoIP lookups, default: false
  geoipOptions:                  # only needed for automatic GeoIP database updates
    accountId: "your maxmind account id"
    licenseKey: "your maxmind license key"
```

> A section left empty (for example `options:` with every line commented out) loses its default values: `healthcheck` becomes `false`. Remove the whole section instead.

## Docker

### Image tags

| Tag | Content |
|---|---|
| `latest` | Latest release |
| `v1` | Latest `1.x.x` release |
| `v1.0.14` | A specific release |
| `develop`, `main`, … | Latest build of a branch |

### Run the container

```
docker run -d \
   --name=emby_exporter \
   -p 9210:9210 \
   -e TZ=Europe/Paris \
   -e CONFIG_FILE=/config/config.yml \
   -v '/path/to/your/config/:/config/' \
   bagul/goemby_exporter:latest
```

`CONFIG_FILE` must be the path of the file **inside** the container.

The image has a Docker `HEALTHCHECK` on `http://localhost:9210/health`. If you set `options.healthcheck: false` or change `exporter.port`, disable it with `--no-healthcheck` (or `healthcheck: disable: true` in docker compose), otherwise the container is reported as `unhealthy`.

### Build the image

The image embeds a MaxMind database: download the GeoLite2 City database from [MaxMind](https://www.maxmind.com/en/home) and place it at the root of the project as `geoip.mmdb`, then:

```
docker build --build-arg VERSION=dev -t emby_exporter .
```

## MaxMind GeoIP Database

The exporter supports GeoIP lookups to provide geographic information about client connections. This requires a MaxMind GeoLite2 City database.

For more detailed information on MaxMind database updates, see their [developer documentation](https://dev.maxmind.com/geoip/updating-databases/#directly-downloading-databases).

### Options for using GeoIP:

1. **Manual database download**:
   - Register for a free account at [MaxMind](https://www.maxmind.com/en/geolite2/signup)
   - Download the GeoLite2 City database (.mmdb file)
   - Place it at the root of the project or in your Docker volume as `geoip.mmdb`
   - Enable GeoIP in your config with `geoip: true`

2. **Automatic database updates**:
   - Register for a free account at [MaxMind](https://www.maxmind.com/en/geolite2/signup)
   - Find your Account ID ([How to find your Account ID](https://support.maxmind.com/hc/en-us/articles/4412951066779-Find-my-Account-ID))
   - Generate or access your License Key ([MaxMind License Keys page](https://www.maxmind.com/en/accounts/current/license-key))
   - Add these credentials to your config file:
   ```yaml
   options:
     geoip: true
     geoipOptions:
       accountId: "your_account_id_here"
       licenseKey: "your_license_key_here"
   ```
   - The exporter will automatically download and update the GeoIP database

### Custom database path

You can store the GeoIP database elsewhere with `GEOIP_DB` (or `-g`). With automatic updates, keep it in a volume so the updated database survives a container restart:

```
docker run -d \
   --name=emby_exporter \
   -p 9210:9210 \
   -e TZ=Europe/Paris \
   -e CONFIG_FILE=/config/config.yml \
   -e GEOIP_DB=/config/GeoLite2-City.mmdb \
   -v '/path/to/your/config/:/config/' \
   bagul/goemby_exporter:latest
```