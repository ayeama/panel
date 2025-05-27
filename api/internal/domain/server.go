package domain

type Server struct {
	Id     string
	Name   string
	Status string

	Container *Container
	Node      *Node
	Manifest  *Manifest
}
