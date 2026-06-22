# transmission-otel

Polls a [Transmission](https://transmissionbt.com/) RPC endpoint and exports torrent/session metrics via OpenTelemetry (OTLP/HTTP).

## Usage

# Docker 🐳

NOTICE: When running in docker make sure you can reach the Transmission RPC endpoint from the container.

## CLI
```sh
docker run -d --name transmission-otel \
  -e TRANSMISSION_HOST=<host> \
  -e TRANSMISSION_PORT=<port> \
  -e TRANSMISSION_USER=<user> \
  -e TRANSMISSION_PASSWORD=<password> \
  -e TRANSMISSION_OTEL_ENDPOINT=<otel_endpoint> \
  git.aads.cloud/aad/bitcoind-metrics-exporter:latest
```

## Compose

```yaml
services:
  transmission-otel:
    image: git.aads.cloud/aad/bitcoind-metrics-exporter:latest
    environment:
      TRANSMISSION_HOST: <host>
      TRANSMISSION_PORT: <port>
      TRANSMISSION_USER: <user>
      TRANSMISSION_PASSWORD: <password>
      TRANSMISSION_OTEL_ENDPOINT: <otel_endpoint>
```


# Install as systemd service
```bash
curl -fsSL https://git.aads.cloud/api/packages/aad/debian/repository.key | sudo gpg --dearmor -o /usr/share/keyrings/transmission-otel.gpg
echo "deb [signed-by=/usr/share/keyrings/transmission-otel.gpg] https://git.aads.cloud/api/packages/aad/debian stable main" | sudo tee /etc/apt/sources.list.d/transmission-otel.list
sudo apt update
sudo apt install transmission-otel
sudo $EDITOR /etc/transmission-otel/transmission-otel.env
sudo systemctl start transmission-otel
```


## Configuration

All settings are read from `TRANSMISSION_`-prefixed environment variables.

| Variable | Default | Description |
|---|---|---|
| `TRANSMISSION_SCHEME` | `http` | Scheme used to reach the Transmission RPC endpoint. |
| `TRANSMISSION_HOST` | `127.0.0.1` | Transmission RPC host. |
| `TRANSMISSION_PORT` | `9091` | Transmission RPC port. |
| `TRANSMISSION_PATH` | `/transmission/rpc` | Transmission RPC path. |
| `TRANSMISSION_USER` | _(empty)_ | RPC basic auth username, if required. |
| `TRANSMISSION_PASSWORD` | _(empty)_ | RPC basic auth password, if required. |
| `TRANSMISSION_OTEL_ENDPOINT` | `http://localhost:4318` | OTLP/HTTP collector endpoint metrics are exported to. Ignored in debug mode. |
| `TRANSMISSION_FETCH_INTERVAL` | `10` | Seconds between polls of Transmission and metric exports. |
| `TRANSMISSION_ERROR_CHECK_INTERVAL` | `60` | Seconds between checks for torrents in an error state. This requires listing every torrent, so it's checked less often than other metrics to stay cheap on servers managing very large numbers of torrents. The reported value holds steady between checks. |
| `TRANSMISSION_SPACE_CHECK_PATH` | _(empty)_ | Comma-separated list of filesystem paths to report free space for (e.g. `/downloads,/media`). Skipped entirely if unset. |
| `TRANSMISSION_DEBUG` | `false` | When `true`, skips the OTEL exporter entirely and just logs fetched metrics to stdout each cycle — useful for checking connectivity without a collector running. |

## Metrics

| Metric | Type | Labels | Description |
|---|---|---|---|
| `torrent_count` | gauge | – | Total number of torrents. |
| `torrent_status_count` | gauge | `status` (`active`, `paused`, `error`) | Number of torrents in each state. |
| `port_open` | gauge | – | Whether the Transmission peer port is open (`1`) or closed (`0`). |
| `download_speed_bytes` | gauge | – | Current download speed, bytes/sec. |
| `upload_speed_bytes` | gauge | – | Current upload speed, bytes/sec. |
| `free_space_bytes` | gauge | `path` | Free disk space at each path configured via `TRANSMISSION_SPACE_CHECK_PATH`. |
| `total_download_bytes` | gauge | – | Total download bytes. |
| `total_upload_bytes` | gauge | – | Total upload bytes. |



No further registration is needed — every declared gauge is wired into the OTel meter automatically.
