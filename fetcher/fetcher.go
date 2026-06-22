package fetcher

import (
	"context"
	"time"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/sirupsen/logrus"

	"git.aads.cloud/aad/transmission-otel/config"

	otelmetrics "git.aads.cloud/aad/transmission-otel/otel/metrics"
)

var log = logrus.WithFields(logrus.Fields{
	"prefix": "fetcher",
})

// Start runs the fetch loop forever, polling Transmission and updating
// metrics every config.C.FetchInterval seconds.
func Start() {
	runner, err := NewRunner()
	if err != nil {
		log.WithError(err).Fatal("Failed to create fetcher runner")
	}

	for {
		runner.run()
		time.Sleep(time.Duration(config.C.FetchInterval) * time.Second)
	}
}

type Runner struct {
	client         *transmissionrpc.Client
	lastErrorCheck time.Time
}

func NewRunner() (*Runner, error) {
	client, err := New()
	if err != nil {
		return nil, err
	}
	return &Runner{client: client}, nil
}

func (r *Runner) run() {
	ctx := context.Background()
	start := time.Now()

	snapshot, err := GetSessionSnapshot(ctx, r.client)
	if err != nil {
		log.WithError(err).Error("Failed to fetch session stats")
		return
	}
	otelmetrics.TorrentCount.Set(snapshot.TorrentCount)
	otelmetrics.TorrentStatus.Set("active", snapshot.ActiveCount)
	otelmetrics.TorrentStatus.Set("paused", snapshot.PausedCount)
	otelmetrics.DownloadSpeed.Set(snapshot.DownloadSpeed)
	otelmetrics.UploadSpeed.Set(snapshot.UploadSpeed)
	otelmetrics.TotalDownload.Set(snapshot.TotalDownload)
	otelmetrics.TotalUpload.Set(snapshot.TotalUpload)

	if time.Since(r.lastErrorCheck) >= time.Duration(config.C.ErrorCheckInterval)*time.Second {
		errorCount, err := GetTorrentErrorCount(ctx, r.client)
		if err != nil {
			log.WithError(err).Error("Failed to fetch torrent error count")
			return
		}
		otelmetrics.TorrentStatus.Set("error", errorCount)
		r.lastErrorCheck = time.Now()
	}

	portStatus, err := PortStatus(ctx, r.client)
	if err != nil {
		log.WithError(err).Error("Failed to fetch port status")
		return
	}
	otelmetrics.PortOpen.SetBool(portStatus)

	if config.C.SpaceCheckPath != "" {
		spaceInfo, err := FreeSpace(ctx, r.client, config.C.SpaceCheckPath)
		if err != nil {
			log.WithError(err).Error("Failed to fetch free space")
			return
		}
		for _, s := range spaceInfo {
			otelmetrics.FreeSpace.Set(s.Path, s.FreeSpace)
		}
	}

	log.WithFields(logrus.Fields{
		"torrent_count": snapshot.TorrentCount,
		"duration":      time.Since(start),
	}).Info("Fetched torrent stats")
}
