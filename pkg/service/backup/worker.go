// Copyright (C) 2017 ScyllaDB

package backup

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/scylladb/go-log"
	"github.com/scylladb/scylla-manager/v3/pkg/metrics"
	"github.com/scylladb/scylla-manager/v3/pkg/scyllaclient"
	. "github.com/scylladb/scylla-manager/v3/pkg/service/backup/backupspec"
	"github.com/scylladb/scylla-manager/v3/pkg/util/parallel"
	"github.com/scylladb/scylla-manager/v3/pkg/util/uuid"
)

// hostInfo groups target host properties needed for backup.
type hostInfo struct {
	DC        string
	IP        string
	ID        string
	Location  Location
	RateLimit DCLimit
}

func (h hostInfo) String() string {
	return h.IP
}

// snapshotDir represents a remote directory containing a table snapshot.
type snapshotDir struct {
	Host     string
	Unit     int64
	Path     string
	Keyspace string
	Table    string
	Version  string
	Progress *RunProgress

	SkippedBytesOffset int64
	NewFilesSize       int64
}

func (sd snapshotDir) String() string {
	return fmt.Sprintf("%s: %s/%s", sd.Host, sd.Keyspace, sd.Table)
}

// workerTools is an intersection of fields and methods
// useful for both worker and restoreWorker.
type workerTools struct {
	ClusterID   uuid.UUID
	ClusterName string
	TaskID      uuid.UUID
	RunID       uuid.UUID
	Config      Config
	Client      *scyllaclient.Client
	Logger      log.Logger
}

// worker is responsible for coordinating backup procedure.
type worker struct {
	workerTools

	Metrics       metrics.BackupM
	SnapshotTag   string
	Units         []Unit
	Schema        *bytes.Buffer
	OnRunProgress func(ctx context.Context, p *RunProgress)
	// ResumeUploadProgress populates upload stats of the provided run progress
	// with previous run progress.
	// If there is no previous run there should be no update.
	// It's required to provide Size as current disk size of files.
	ResumeUploadProgress func(ctx context.Context, p *RunProgress)

	// Cache for host snapshotDirs
	snapshotDirs map[string][]snapshotDir
	mu           sync.Mutex

	memoryPool *sync.Pool
}

func (w *worker) WithLogger(logger log.Logger) *worker {
	w.Logger = logger
	return w
}

func (w *worker) hostSnapshotDirs(h hostInfo) []snapshotDir {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.snapshotDirs[h.IP]
}

func (w *worker) setSnapshotDirs(h hostInfo, dirs []snapshotDir) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.snapshotDirs == nil {
		w.snapshotDirs = make(map[string][]snapshotDir)
	}

	w.snapshotDirs[h.IP] = dirs
}

// cleanup resets global stats for each agent.
func (w *worker) cleanup(ctx context.Context, hi []hostInfo) {
	if err := hostsInParallel(hi, parallel.NoLimit, func(h hostInfo) error {
		return w.Client.RcloneResetStats(context.Background(), h.IP)
	}); err != nil {
		w.Logger.Error(ctx, "Failed to reset stats", "error", err)
	}
}
