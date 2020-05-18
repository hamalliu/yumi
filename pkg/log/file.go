package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const logdir = "logfile"
const dateformat = "2006-01-02"
const timeformat = "2006-01-02 15-03-04.999"

type File struct {
	mux sync.Mutex
	f   *os.File

	prefix  string
	maxSize int64
	day     bool
}

func New(prefix string, maxsize int64, day bool) *File {
	_ = os.Mkdir(logdir, 0644)
	return &File{
		prefix:  prefix,
		maxSize: maxsize,
		day:     day,
	}
}

func (m *File) Write(b []byte) (int, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if err := m.initFile(); err != nil {
		return 0, err
	}

	return m.f.Write(b)
}

func (m *File) initFile() error {
	var (
		curFileName = fmt.Sprintf("%s/%s%s", logdir, m.prefix, time.Now().Format(dateformat))
		err         error
	)
	if m.day && (m.f == nil || m.f.Name() != curFileName) {
		if m.f, err = os.OpenFile(curFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
			return err
		}
	}

	if m.maxSize != 0 {
		if finfo, err := m.f.Stat(); err != nil {
			_ = m.f.Close()
			if m.f, err = os.OpenFile(curFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
				return err
			}

			finfo, _ = m.f.Stat()
			if m.maxSize < finfo.Size() {
				_ = m.f.Close()
				_ = os.Rename(curFileName, fmt.Sprintf("%s/%s%s", logdir, m.prefix, time.Now().Format(timeformat)))
				if m.f, err = os.OpenFile(curFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
					return err
				}
			}
		} else {
			if m.maxSize < finfo.Size() {
				_ = m.f.Close()
				_ = os.Rename(curFileName, fmt.Sprintf("%s/%s%s", logdir, m.prefix, time.Now().Format(timeformat)))
				if m.f, err = os.OpenFile(curFileName, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
