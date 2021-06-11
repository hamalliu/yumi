package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// LogOption ...
type LogOption struct {
	F func(*LogOptions)
}

// LogOptions ...
type LogOptions struct {
	StorageDir   string
	FileName     string
	IsOutputStd  bool
	FileMaxAge   time.Duration
	RotationTime time.Duration
}

func (lo *LogOptions) NewFileWriter(subDir string) (io.Writer, error) {
	storageDir := ""
	if subDir != "" {
		storageDir = filepath.Join(lo.StorageDir, subDir)
	}
	r, err := rotatelogs.New(storageDir,
		rotatelogs.WithLinkName(filepath.Join(storageDir, fmt.Sprintf("%s.log", lo.FileName))), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(lo.FileMaxAge),                                                      // 文件最大保存时间
		rotatelogs.WithRotationTime(lo.RotationTime),                                              // 日志切割时间间隔
	)
	if err != nil {
		return nil, err
	}

	if lo.IsOutputStd {
		return io.MultiWriter(os.Stdout, r), nil
	}

	return r, nil
}

// SetStorageDir ...
func SetStorageDir(storageDir string) LogOption {
	return LogOption{
		F: func(lo *LogOptions) {
			lo.StorageDir = storageDir
		},
	}
}

// SetFileName ...
func SetFileName(fileName string) LogOption {
	return LogOption{
		F: func(lo *LogOptions) {
			lo.FileName = fileName
		},
	}
}

// SetIsOutputStd ...
func SetIsOutputStd(isOutputStd bool) LogOption {
	return LogOption{
		F: func(lo *LogOptions) {
			lo.IsOutputStd = isOutputStd
		},
	}
}

// SetFileMaxAge ...
func SetFileMaxAge(fileMaxAge time.Duration) LogOption {
	return LogOption{
		F: func(lo *LogOptions) {
			lo.FileMaxAge = fileMaxAge
		},
	}
}

// SetRotationTime ...
func SetRotationTime(rotationTime time.Duration) LogOption {
	return LogOption{
		F: func(lo *LogOptions) {
			lo.RotationTime = rotationTime
		},
	}
}
