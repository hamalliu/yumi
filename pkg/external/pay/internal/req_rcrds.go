package internal

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

//ReqRcrd ...
type ReqRcrd struct {
	f *File
}

//NewReqRcrd ...
func NewReqRcrd(name string) *ReqRcrd {
	f := newFile(name, defaultMemory, true)

	return &ReqRcrd{f}
}

//AddRcrds ...
func (rr *ReqRcrd) AddRcrds(method, reqURL string, reqHeader http.Header, reqBody []byte, respHeader http.Header, respBody []byte) {
	_, _ = rr.f.WriteString(time.Now().Format("2006-01-02 15:04:05.999: ") + method + " reqUrl: ")
	_, _ = rr.f.WriteString(reqURL)
	_, _ = rr.f.WriteString("\n")

	_, _ = rr.f.WriteString("reqHeader:\n")
	for i := range reqHeader {
		_, _ = rr.f.WriteString(reqHeader.Get(i))
		_, _ = rr.f.WriteString("\n")
	}

	_, _ = rr.f.WriteString("reqBody:\n")
	_, _ = rr.f.Write(reqBody)
	_, _ = rr.f.WriteString("\n")

	_, _ = rr.f.WriteString("respHeader:\n")
	for i := range respHeader {
		_, _ = rr.f.WriteString(respHeader.Get(i))
		_, _ = rr.f.WriteString("\n")
	}

	_, _ = rr.f.WriteString("respBody:\n")
	_, _ = rr.f.Write(respBody)
	_, _ = rr.f.WriteString("\n")
	_, _ = rr.f.WriteString("========================================================================================================================")
	_, _ = rr.f.WriteString("\n")

	_ = rr.f.Sync()
}

const logDir = "pay_records"
const dateFormat = "2006-01-02"
const timeFormat = "2006-01-02 15-03-04.999"
const defaultMemory = 32 << 20 // 32 MB

//File ...
type File struct {
	mux sync.Mutex
	f   *os.File

	prefix  string
	maxSize int64
	day     bool
}

func newFile(prefix string, maxsize int64, day bool) *File {
	_ = os.Mkdir(logDir, 0644)
	return &File{
		prefix:  prefix,
		maxSize: maxsize,
		day:     day,
	}
}

func (m *File) Write(b []byte) (int, error) {
	var (
		err error
	)

	m.mux.Lock()
	defer m.mux.Unlock()

	if err = m.initFile(); err != nil {
		return 0, err
	}

	return m.f.Write(b)
}

//WriteString ...
func (m *File) WriteString(str string) (int, error) {
	var (
		err error
	)

	m.mux.Lock()
	defer m.mux.Unlock()

	if err = m.initFile(); err != nil {
		return 0, err
	}

	return m.f.WriteString(str)
}

//Sync ...
func (m *File) Sync() error {
	return m.f.Sync()
}

func (m *File) initFile() error {
	var (
		curFileName = fmt.Sprintf("%s/%s%s", logDir, m.prefix, time.Now().Format(dateFormat))
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
				_ = os.Rename(curFileName, fmt.Sprintf("%s/%s%s", logDir, m.prefix, time.Now().Format(timeFormat)))
				if m.f, err = os.OpenFile(curFileName, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
					return err
				}
			}
		} else {
			if m.maxSize < finfo.Size() {
				_ = m.f.Close()
				_ = os.Rename(curFileName, fmt.Sprintf("%s/%s%s", logDir, m.prefix, time.Now().Format(timeFormat)))
				if m.f, err = os.OpenFile(curFileName, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
