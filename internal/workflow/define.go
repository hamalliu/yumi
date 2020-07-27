package workflow

import (
	"time"
)

type Define struct {
	Name       string
	CreateDate time.Time
	Creator    string
}

type Node struct {
	Name string
}
