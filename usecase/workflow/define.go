package workflow

import (
	"time"
)

//Define ...
type Define struct {
	Name       string
	CreateDate time.Time
	Creator    string
}

//Node ...
type Node struct {
	Name string
}
