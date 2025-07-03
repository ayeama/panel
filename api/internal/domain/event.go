package domain

type RuntimeEvent interface{}

type RuntimeEventServerStatusChanged struct {
	ServerId    string
	ContainerId string
}
