package runtime

import (
	"archive/zip"
	"io"

	"github.com/ayeama/panel/internal/types"
)

type Runtime interface {
	ImageReadMany() ([]types.Image, error)

	InstanceCreate(imageID string, resources types.InstanceResources, webhooks []string) (types.Instance, error)
	InstanceRead(id string) (types.Instance, error)
	InstanceReadMany() ([]types.Instance, error)
	InstanceDelete(id string) error
	InstanceStart(id string) error
	InstanceStop(id string) error
	InstanceAttach(id string, stdin io.Reader, stdout io.Writer, stderr io.Writer, ready chan bool) error
	InstanceStats(id string, stats chan types.InstanceStat) error
	InstanceLogs(id string, logs chan string) error
	InstanceBackup(id string, manifest *types.InstanceBackupManifest, zw *zip.Writer) error
	InstanceRestore(id string, manifest *types.InstanceBackupManifest, zr *zip.Reader) error

	Events(events chan types.Event, cancel chan bool) error
}
