package otelmetrics

var PortOpen = NewGauge("port_open", "Whether the Transmission peer port is open (1) or closed (0)")
var TorrentCount = NewGauge("torrent_count", "Total number of torrents")
var FreeSpace = NewLabeledGauge("free_space_bytes", "Free disk space in bytes", "path")
var DownloadSpeed = NewGauge("download_speed_bytes", "Current download speed in bytes/sec")
var UploadSpeed = NewGauge("upload_speed_bytes", "Current upload speed in bytes/sec")
var TorrentStatus = NewLabeledGauge("torrent_status_count", "Number of torrents per status", "status")
var TotalDownload = NewGauge("total_download_bytes", "Total downloaded bytes")
var TotalUpload = NewGauge("total_upload_bytes", "Total uploaded bytes")
