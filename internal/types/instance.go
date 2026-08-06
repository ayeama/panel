package types

type Instance struct {
	ID        string
	Name      string
	Image     string
	Status    string
	Ports     map[string]string
	Resources InstanceResources
	Webhooks  []string
}

type InstanceCreate struct {
	Image     string
	Resources InstanceResources
	Webhooks  []string
}

type InstanceResources struct {
	CPU    float64
	Memory float64
	Disk   float64
}

type InstanceStat struct {
	CPUPercent     float64
	MemoryPercent  float64
	DiskPercent    float64
	NetworkTXBytes uint64
	NetworkRXBytes uint64
}

type InstanceBackupManifest struct {
	Instance

	Version string
	Mounts  []InstanceBackupManifestMount
}

type InstanceBackupManifestMount struct {
	ID          string
	Destination string
}
