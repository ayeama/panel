package domain

import "time"

type Session struct {
	Id      string
	UserId  string
	Expires time.Time
}
