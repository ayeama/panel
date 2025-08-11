package domain

type Container struct {
	Id     string
	Name   string
	Status string
	Ports  []string
}

type ContainerStat struct {
	Cpu    float64
	Memory float64
}
