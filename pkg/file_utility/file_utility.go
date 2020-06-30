package file_utility

import (
	"fmt"
	"io"
	"os"
)

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

func CreateDir(path string) error {
	if !ExistDir(path) {
		if err := os.MkdirAll(path, 0644); err != nil {
			return err
		}
	}
	return nil
}

func DeleteDir(path string) error {
	if ExistDir(path) {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func CopyFile(srcPath, destPath string) error {
	sf, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer sf.Close()

	if ExistFile(destPath) {
		return fmt.Errorf("file already exists")
	}

	df, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}

	return nil
}
