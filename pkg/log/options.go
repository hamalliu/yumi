package log

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// Option ...
type Option struct {
	F func(*Options)
}

// Options ...
type Options struct {
	StorageDir    string
	FileName      string
	IsOutputStd   bool
	RotationCount uint
	RotationTime  time.Duration
}

// NewFileOutput ...
func (lo *Options) NewFileOutput(level Level) (io.Writer, error) {
	fileOutput, err := lo.newFileOutput(level.ToString())
	if err != nil {
		return nil, err
	}
	return fileOutput, nil
}

func (lo *Options) defaultSet() {
	lo.StorageDir = "logfile"
	lo.FileName = "yumi"
	lo.RotationCount = 7
	lo.RotationTime = time.Hour * 24
	lo.IsOutputStd = true
}

func (lo *Options) newFileOutput(subDir string) (io.Writer, error) {
	lo.defaultSet()

	storageDir := lo.StorageDir
	if subDir != "" {
		storageDir = filepath.Join(lo.StorageDir, subDir)
	}
	r, err := rotatelogs.New(filepath.Join(storageDir, lo.FileName+"-%Y%m%d.log"),
		rotatelogs.WithLinkName(filepath.Join(storageDir, fmt.Sprintf("%s.log", lo.FileName))), // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(lo.RotationCount),                                         // 文件最多保存多少个
		rotatelogs.WithRotationTime(lo.RotationTime),                                           // 轮询日志切割时间间隔
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// SetStorageDir ...
func SetStorageDir(storageDir string) Option {
	return Option{
		F: func(lo *Options) {
			lo.StorageDir = storageDir
		},
	}
}

// SetFileName ...
func SetFileName(fileName string) Option {
	return Option{
		F: func(lo *Options) {
			lo.FileName = fileName
		},
	}
}

// SetIsOutputStd ...
func SetIsOutputStd(isOutputStd bool) Option {
	return Option{
		F: func(lo *Options) {
			lo.IsOutputStd = isOutputStd
		},
	}
}

// SetRotationCount ...
func SetRotationCount(n uint) Option {
	return Option{
		F: func(lo *Options) {
			lo.RotationCount = n
		},
	}
}

// SetRotationTime ...
func SetRotationTime(rotationTime time.Duration) Option {
	return Option{
		F: func(lo *Options) {
			lo.RotationTime = rotationTime
		},
	}
}
