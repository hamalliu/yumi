package log

import (
	"fmt"
	"io"
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
	StorageDir    string
	FileName      string
	IsOutputStd   bool
	RotationCount uint
	RotationTime  time.Duration
}

func (lo *LogOptions) NewFileOutput(upOutput io.Writer, level Level) (io.Writer, error) {
	fileOutput, err := lo.newFileOutput(level.ToString())
	if err != nil {
		return nil, err
	}
	if upOutput != nil {
		return io.MultiWriter(fileOutput, upOutput), nil
	}

	return fileOutput, nil
}

func (lo *LogOptions) defaultSet() {
	if lo.StorageDir == "" {
		lo.StorageDir = "logfile"
	}
	if lo.FileName == "" {
		lo.FileName = "yumi"
	}
	if lo.RotationCount == 0 {
		lo.RotationCount = 30
	}
	if lo.RotationTime == 0 {
		lo.RotationTime = time.Minute
	}
}

func (lo *LogOptions) newFileOutput(subDir string) (io.Writer, error) {
	lo.defaultSet()

	storageDir := ""
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
func SetRotationCount(n uint) LogOption {
	return LogOption{
		F: func(lo *LogOptions) {
			lo.RotationCount = n
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
