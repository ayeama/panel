package domain

type Sidecar struct {
	Id          string
	ContainerId string
	ServerId    string

	Container *Container
}
