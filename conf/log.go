package conf

import "yumi/pkg/log"

// Log ...
type Log struct {
	Level       log.Level
	StorageDir  string
	FileName    string
	IsOutputStd bool
}

// Options ...
func (l *Log) Options() []log.Option {
	opts := []log.Option{}
	opts = append(opts, log.SetFileName(l.FileName))
	opts = append(opts, log.SetIsOutputStd(l.IsOutputStd))
	opts = append(opts, log.SetStorageDir(l.StorageDir))

	return opts
}
