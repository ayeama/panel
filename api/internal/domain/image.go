package domain

import (
	"path"
)

type Image struct {
	Id  string
	Tag string

	Variables *map[string]string
}

func (i *Image) String() string {
	return path.Base(i.Tag)
}
