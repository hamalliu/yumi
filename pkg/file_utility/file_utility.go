package file_utility

import "os"

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

func CreateDir(path string) bool {
	if !ExistDir(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return false
		}
	}
	return true
}

func DeleteDir(path string) bool {
	if ExistDir(path) {
		if err := os.RemoveAll(path); err != nil {
			return false
		}
	}
	return true
}
