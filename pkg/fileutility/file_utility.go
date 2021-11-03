package fileutility

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

//ExistDir ...
func ExistDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	} else if f.IsDir() {
		return true
	} else {
		return false
	}
}

//ExistFile ...
func ExistFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	} else if f.IsDir() {
		return false
	} else {
		return true
	}
}

//CreateDir ...
func CreateDir(path string) error {
	if !ExistDir(path) {
		if err := os.MkdirAll(path, 0644); err != nil {
			return err
		}
	}
	return nil
}

//DeleteDir ...
func DeleteDir(path string) error {
	if ExistDir(path) {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

//CopyFile ...
func CopyFile(srcPath, destPath string) error {
	sf, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() { _ = sf.Close() }()

	if ExistFile(destPath) {
		return fmt.Errorf("file already exists")
	}

	df, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() { _ = df.Close() }()

	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}

	return nil
}

//WriteFile ...
func WriteFile(src io.Reader, destPath string) error {
	if ExistFile(destPath) {
		return fmt.Errorf("file already exists")
	}

	df, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() { _ = df.Close() }()

	_, err = io.Copy(df, src)
	if err != nil {
		return err
	}

	return nil
}

//GetFileExtension ...
func GetFileExtension(fileName string, withoutDot bool) string {
	if withoutDot {
		return fileName[strings.LastIndex(fileName, ".")+1:]
	}

	return fileName[strings.LastIndex(fileName, "."):]
}

//GetModTime ...
func GetModTime(path string) time.Time {
	s, _ := os.Stat(path)
	return s.ModTime()
}
