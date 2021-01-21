package docmanager

import (
	"fmt"
	"os"
	"path"
	"strings"

	"yumi/conf"
	"yumi/pkg/fileutility"
	"yumi/pkg/types"
)

const (
	//FileTypeText ...
	FileTypeText = "text"
	//FileTypeSpreadsheet ...
	FileTypeSpreadsheet = "spreadsheet"
	//FileTypePresentation ...
	FileTypePresentation = "presentation"
)

const (
	//InteralFileExtDocx ...
	InteralFileExtDocx = ".docx"
	//InteralFileExtXlsx ...
	InteralFileExtXlsx = ".xlsx"
	//InteralFileExtPptx ...
	InteralFileExtPptx = ".pptx"
)

var documentExts = types.ArrayString{
	".doc", ".docx", ".docm", ".dot", ".dotx", ".dotm", ".odt", ".fodt", ".ott", ".rtf", ".txt", ".html", ".htm",
	".mht", ".pdf", ".djvu", ".fb2", ".epub", ".xps"}

var spreadsheetExts = types.ArrayString{
	".xls", ".xlsx", ".xlsm", ".xlt", ".xltx", ".xltm", ".ods", ".fods", ".ots", ".csv"}

var presentationExts = types.ArrayString{
	".pps", ".ppsx", ".ppsm", ".ppt", ".pptx", ".pptm", ".pot", ".potx", ".potm", ".odp", ".fodp", ".otp"}

//FileUtility ...
type FileUtility struct {
	cfg conf.Document
}

//GetFileName ...
func (fu FileUtility) GetFileName(fileName string, withoutExtension bool) string {
	if fileName == "" {
		return ""
	}

	fileName = strings.ToLower(fileName)
	s := strings.LastIndex(fileName, "/")
	if s != -1 {
		fileName = fileName[s+1:]
	}

	if withoutExtension {
		fileName = strings.TrimSuffix(fileName, fu.GetFileExtension(fileName, false))
	}

	return fileName
}

//GetFileExtension ...
func (fu FileUtility) GetFileExtension(fileName string, withoutDot bool) string {
	s := strings.LastIndex(fileName, ".")
	if s == -1 {
		return ""
	}

	if withoutDot {
		return fileName[s+1:]
	}

	return fileName[s:]
}

//GetFileType ...
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

//StoragePath ...
func (fu FileUtility) StoragePath(fileName, userID string) string {
	directory := path.Join(fu.cfg.StoragePath, userID)
	_ = fileutility.CreateDir(directory)

	fileName = fu.GetFileName(fileName, false)

	return path.Join(directory, fileName)
}

//ForcesavePath ...
func (fu FileUtility) ForcesavePath(fileName, userID string, create bool) string {
	directory := path.Join(fu.cfg.StoragePath, userID)
	if !fileutility.ExistDir(directory) {
		return ""
	}

	directory = path.Join(directory, fileName+"-history")
	if !create && !fileutility.ExistDir(directory) {
		return ""
	}

	_ = fileutility.CreateDir(directory)
	directory = path.Join(directory, fileName)
	if !create && !fileutility.ExistDir(directory) {
		return ""
	}

	return directory
}

//HistoryPath ...
func (fu FileUtility) HistoryPath(fileName, userID string, create bool) string {
	directory := path.Join(fu.cfg.StoragePath, userID)
	if !fileutility.ExistDir(directory) {
		return ""
	}

	directory = path.Join(directory, fileName+"-history")
	if !create && !fileutility.ExistDir(path.Join(directory, "1")) {
		return ""
	}

	return directory
}

//VersionPath ...
func (fu FileUtility) VersionPath(fileName, userID string, version int) string {
	directory := fu.HistoryPath(fileName, userID, true)

	return path.Join(directory, fmt.Sprintf("%d", version))
}

//PrevFilePath ...
func (fu FileUtility) PrevFilePath(fileName, userID string, version int) string {
	directory := fu.VersionPath(fileName, userID, version)

	return path.Join(directory, "prev"+fu.GetFileExtension(fileName, false))
}

//DiffPath ...
func (fu FileUtility) DiffPath(fileName, userID string, version int) string {
	directory := fu.VersionPath(fileName, userID, version)

	return path.Join(directory, "diff.zip")
}

//ChangesPath ...
func (fu FileUtility) ChangesPath(fileName, userID string, version int) string {
	directory := fu.VersionPath(fileName, userID, version)

	return path.Join(directory, "changes.txt")
}

//KeyPath ...
func (fu FileUtility) KeyPath(fileName, userID string, version int) string {
	directory := fu.VersionPath(fileName, userID, version)

	return path.Join(directory, "key.txt")
}

//GetCorrectName ...
func (fu FileUtility) GetCorrectName(fileName, userID string) string {
	baseName := fu.GetFileName(fileName, true)
	ext := fu.GetFileExtension(fileName, false)

	name := baseName + ext
	index := 1

	for {
		if fileutility.ExistFile(fu.StoragePath(fileName, userID)) {
			name = fmt.Sprintf("%s(%d)%s", baseName, index, ext)
			index++
		}
		break
	}

	return name
}

//GetInternalExtension ...
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

//CleanFolderRecursive ...
func (fu FileUtility) CleanFolderRecursive(floder string, me bool) error {
	if err := os.RemoveAll(floder); err != nil {
		return err
	}

	if !me {
		return os.Mkdir(floder, 0644)
	}

	return nil
}

//AllowUploadExtension ...
func (fu FileUtility) AllowUploadExtension(ext string) bool {
	if fu.cfg.EditedDocs.IndexOf(ext) != -1 ||
		fu.cfg.ViewedDocs.IndexOf(ext) != -1 ||
		fu.cfg.ConvertedDocs.IndexOf(ext) != -1 {
		return true
	}

	return false
}
