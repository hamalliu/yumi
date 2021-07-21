package conf

import "yumi/pkg/log"

type Log struct {
	Level       log.Level
	StorageDir  string
	FileName    string
	IsOutputStd bool
}

// Options ...
func (l *Log) Options() []log.LogOption {
	opts := []log.LogOption{}
	opts = append(opts, log.SetFileName(l.FileName))
	opts = append(opts, log.SetIsOutputStd(l.IsOutputStd))
	opts = append(opts, log.SetStorageDir(l.StorageDir))

	return opts
}
