# transmission-otel

Polls a [Transmission](https://transmissionbt.com/) RPC endpoint and exports torrent/session metrics via OpenTelemetry (OTLP/HTTP).

## Usage

```sh
go run .
```

By default it connects to `http://127.0.0.1:9091/transmission/rpc` and exports metrics to `http://localhost:4318` every 10s. Configure via environment variables (see below).

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

### Adding a new metric

Declare a gauge in `otel/metrics/metrics.go`:

```go
var MyMetric = NewGauge("my_metric", "Description")
// or, for one value per item (e.g. per path/disk):
var MyMetric = NewLabeledGauge("my_metric", "Description", "label_key")
```

Then set its value from the fetch loop in `fetcher/fetcher.go`:

```go
otelmetrics.MyMetric.Set(value)
// or
otelmetrics.MyMetric.Set(label, value)
```

No further registration is needed — every declared gauge is wired into the OTel meter automatically.
