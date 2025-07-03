package domain

type Container struct {
	Id   string
	Port string
}

type ContainerStat struct {
	Cpu    float64
	Memory float64
}
