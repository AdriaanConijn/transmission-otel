package fetcher

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/url"
	"strings"

	"github.com/hekmon/transmissionrpc/v3"

	"git.aads.cloud/aad/transmission-otel/config"
)

func New() (*transmissionrpc.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return NewFromConfig(cfg)
}

func NewFromConfig(cfg config.Config) (*transmissionrpc.Client, error) {
	endpoint := &url.URL{
		Scheme: cfg.Scheme,
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   cfg.Path,
	}
	if cfg.User != "" {
		endpoint.User = url.UserPassword(cfg.User, cfg.Password)
	}

	client, err := transmissionrpc.New(endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating transmissionrpc client: %w", err)
	}
	return client, nil
}

type SessionSnapshot struct {
	TorrentCount  int64
	ActiveCount   int64
	PausedCount   int64
	DownloadSpeed int64
	UploadSpeed   int64
	TotalDownload int64
	TotalUpload   int64
}

func GetSessionSnapshot(ctx context.Context, client *transmissionrpc.Client) (SessionSnapshot, error) {
	stats, err := client.SessionStats(ctx)
	if err != nil {
		return SessionSnapshot{}, fmt.Errorf("session stats: %w", err)
	}
	return SessionSnapshot{
		TorrentCount:  stats.TorrentCount,
		ActiveCount:   stats.ActiveTorrentCount,
		PausedCount:   stats.PausedTorrentCount,
		DownloadSpeed: stats.DownloadSpeed,
		UploadSpeed:   stats.UploadSpeed,
		TotalDownload: stats.CumulativeStats.DownloadedBytes,
		TotalUpload:   stats.CumulativeStats.UploadedBytes,
	}, nil
}

func PortStatus(ctx context.Context, client *transmissionrpc.Client) (bool, error) {
	status, err := client.PortTest(ctx)
	if err != nil {
		return false, fmt.Errorf("port test: %w", err)
	}
	return status, nil
}

func GetTorrentErrorCount(ctx context.Context, client *transmissionrpc.Client) (int64, error) {
	torrents, err := client.TorrentGet(ctx, []string{"id", "error"}, nil)
	if err != nil {
		return 0, fmt.Errorf("listing torrents: %w", err)
	}
	var count int64
	for _, t := range torrents {
		if t.Error != nil && *t.Error != 0 {
			count++
		}
	}
	return count, nil
}

type SpaceInfo struct {
	Path      string
	FreeSpace int64
}

func FreeSpace(ctx context.Context, client *transmissionrpc.Client, paths string) ([]SpaceInfo, error) {
	r := csv.NewReader(strings.NewReader(paths))
	record, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("parsing paths: %w", err)
	}

	results := make([]SpaceInfo, 0, len(record))
	for _, path := range record {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		freeSpace, _, err := client.FreeSpace(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("free space for %q: %w", path, err)
		}
		results = append(results, SpaceInfo{
			Path:      path,
			FreeSpace: int64(freeSpace.Byte()),
		})
	}
	return results, nil
}
