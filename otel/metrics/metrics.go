package otelmetrics

var PortOpen = NewGauge("transmission_port_open", "Whether the Transmission peer port is open (1) or closed (0)")
var TorrentCount = NewGauge("transmission_torrent_count", "Total number of torrents")
var FreeSpace = NewLabeledGauge("transmission_free_space_bytes", "Free disk space in bytes", "path")
var DownloadSpeed = NewGauge("transmission_download_speed_bytes", "Current download speed in bytes/sec")
var UploadSpeed = NewGauge("transmission_upload_speed_bytes", "Current upload speed in bytes/sec")
var TorrentStatus = NewLabeledGauge("transmission_torrent_status_count", "Number of torrents per status", "status")
var TotalDownload = NewCounter("transmission_download_bytes_total", "Total downloaded bytes")
var TotalUpload = NewCounter("transmission_upload_bytes_total", "Total uploaded bytes")
