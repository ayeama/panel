package runtime

import (
	"io"

	"github.com/ayeama/panel/internal/types"
)

type Runtime interface {
	ImageRead(id string) (types.Image, error)
	ImageReadMany() ([]types.Image, error)

	InstanceCreate(imageID string, resources types.InstanceResources) (types.Instance, error)
	InstanceRead(id string) (types.Instance, error)
	InstanceReadMany() ([]types.Instance, error)
	InstanceDelete(id string) error
	InstanceStart(id string) error
	InstanceStop(id string) error
	InstanceAttach(id string, stdin io.Reader, stdout io.Writer, stderr io.Writer, ready chan bool) error
	InstanceStats(id string, stats chan types.InstanceStat) error
	InstanceLogs(id string, logs chan string) error

	Events(events chan types.Event, cancel chan bool) error
}
