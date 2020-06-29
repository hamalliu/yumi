package helplers

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"yumi/pkg/conf"
	"yumi/pkg/file_utility"
	"yumi/pkg/types"
)

const (
	FileTypeText         = "text"
	FileTypeSpreadsheet  = "spreadsheet"
	FileTypePresentation = "presentation"
)

const (
	InteralFileExtDocx = ".docx"
	InteralFileExtXlsx = ".xlsx"
	InteralFileExtPptx = ".pptx"
)

var documentExts = types.ArrayString{
	".doc", ".docx", ".docm", ".dot", ".dotx", ".dotm", ".odt", ".fodt", ".ott", ".rtf", ".txt", ".html", ".htm",
	".mht", ".pdf", ".djvu", ".fb2", ".epub", ".xps"}

var spreadsheetExts = types.ArrayString{
	".xls", ".xlsx", ".xlsm", ".xlt", ".xltx", ".xltm", ".ods", ".fods", ".ots", ".csv"}

var presentationExts = types.ArrayString{
	".pps", ".ppsx", ".ppsm", ".ppt", ".pptx", ".pptm", ".pot", ".potx", ".potm", ".odp", ".fodp", ".otp"}

type FileUtility struct{}

var fileUtility FileUtility

func (fu FileUtility) GetFileName(urlStr string, withoutExtension bool) string {
	if urlStr == "" {
		return ""
	}

	fileName := ""

	cfgOffice := conf.Get().Office
	cfgOffice.TempStorageUrl = cfgOffice.SiteUrl + cfgOffice.TempStorageUrl

	if cfgOffice.TempStorageUrl != "" && strings.Index(urlStr, cfgOffice.TempStorageUrl) == 0 {
		urlObj, err := url.Parse(urlStr)
		if err != nil {
			return ""
		}
		fileName = urlObj.Query().Get("filename")
	} else {
		urlStr = strings.ToLower(urlStr)
		s := strings.LastIndex(urlStr, "/")
		if s == -1 {
			fileName = urlStr
		} else {
			fileName = urlStr[s+1:]
		}
	}

	if withoutExtension {
		fileName = strings.TrimSuffix(fileName, fu.GetFileExtension(fileName, false))
	}

	return fileName
}

func (fu FileUtility) GetFileExtension(fileName string, withoutDot bool) string {
	s := strings.LastIndex(fileName, ".")
	if s == -1 {
		return ""
	}

	if withoutDot {
		return fileName[s+1:]
	} else {
		return fileName[s:]
	}
}

func (fu FileUtility) GetFileType(fileName string) string {
	ext := fu.GetFileExtension(fileName, false)
	switch {
	case documentExts.IndexOf(ext) != -1:
		return FileTypeText
	case spreadsheetExts.IndexOf(ext) != -1:
		return FileTypeSpreadsheet
	case presentationExts.IndexOf(ext) != -1:
		return FileTypePresentation
	default:
		return FileTypeText
	}
}

func (fu FileUtility) CreateDirectory(name string) {
	if _, err := os.Stat(name); err != nil {
		if os.IsExist(err) {
			return
		} else {
			_ = os.Mkdir(name, 0644)
		}
	} else {
		return
	}
}

func (fu FileUtility) StoragePath(fileName, userid string) string {
	cfgOffice := conf.Get().Office

	directory := path.Join(cfgOffice.StoragePath, userid)
	file_utility.CreateDir(directory)

	fileName = fu.GetFileName(fileName, false)

	return path.Join(directory, fileName)
}

func (fu FileUtility) ForcesavePath(fileName, userid string, create bool) string {
	cfgOffice := conf.Get().Office

	directory := path.Join(cfgOffice.StoragePath, userid)
	if !file_utility.ExistDir(directory) {
		return ""
	}

	directory = path.Join(directory, fileName+"-history")
	if !create && !file_utility.ExistDir(directory) {
		return ""
	}

	file_utility.CreateDir(directory)
	directory = path.Join(directory, fileName)
	if !create && !file_utility.ExistDir(directory) {
		return ""
	}

	return directory
}

func (fu FileUtility) HistoryPath(fileName, userid string, create bool) string {
	cfgOffice := conf.Get().Office

	directory := path.Join(cfgOffice.StoragePath, userid)
	if !file_utility.ExistDir(directory) {
		return ""
	}

	directory = path.Join(directory, fileName+"-history")
	if !create && !file_utility.ExistDir(path.Join(directory, "1")) {
		return ""
	}

	return directory
}

func (fu FileUtility) VersionPath(fileName, userid string, version int) string {
	directory := fu.HistoryPath(fileName, userid, true)

	return path.Join(directory, fmt.Sprintf("%d", version))
}

func (fu FileUtility) PrevFilePath(fileName, userid string, version int) string {
	directory := fu.VersionPath(fileName, userid, version)

	return path.Join(directory, "prev"+fu.GetFileExtension(fileName, false))
}

func (fu FileUtility) DiffPath(fileName, userid string, version int) string {
	directory := fu.VersionPath(fileName, userid, version)

	return path.Join(directory, "diff.zip")
}

func (fu FileUtility) ChangesPath(fileName, userid string, version int) string {
	directory := fu.VersionPath(fileName, userid, version)

	return path.Join(directory, "changes.txt")
}

func (fu FileUtility) KeyPath(fileName, userid string, version int) string {
	directory := fu.VersionPath(fileName, userid, version)

	return path.Join(directory, "key.txt")
}

func (fu FileUtility) ChangesUser(fileName, userid string, version int) string {
	directory := fu.VersionPath(fileName, userid, version)

	return path.Join(directory, "user.txt")
}

func (fu FileUtility) GetCorrectName(fileName, userid string) string {
	baseName := fu.GetFileName(fileName, true)
	ext := fu.GetFileExtension(fileName, false)

	name := baseName + ext
	index := 1

	for {
		if file_utility.ExistFile(fu.StoragePath(fileName, userid)) {
			name = fmt.Sprintf("%s(%d)%s", baseName, index, ext)
			index++
		}
		break
	}

	return name
}

func (fu FileUtility) GetInternalExtension(fileType string) string {
	switch {
	case FileTypeText == fileType:
		return InteralFileExtDocx
	case FileTypeSpreadsheet == fileType:
		return InteralFileExtXlsx
	case FileTypePresentation == fileType:
		return InteralFileExtPptx
	default:
		return InteralFileExtDocx
	}
}

func (fu FileUtility) CleanFolderRecursive(floder string, me bool) {

}
