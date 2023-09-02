package filereader

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"stock-data-processing/utils"
)

type FileReaderCollections interface {
}

type FileReader struct {
	FolderPath        string
	SuccessFolderPath string
	FailFolderPath    string
	FileList          []os.DirEntry
	allowedExtFile    []string
}

// NewFileReader initiate new file reader.
// extFile: fill it to specify which file to read by file extension. ex: []string{"txt", "json"}
func NewFileReader(sourceRawFolderPath, successFolderPath, failFolderPath string, extFile []string) (*FileReader, error) {
	// check directory exist
	errmsg := errors.New("directory doesn't exist")
	if _, err := os.Stat(sourceRawFolderPath); os.IsNotExist(err) {
		log.Err(err).Msg(errmsg.Error())
		return nil, err
	}
	if _, err := os.Stat(successFolderPath); os.IsNotExist(err) {
		log.Err(err).Msg(errmsg.Error())
		return nil, err
	}
	if _, err := os.Stat(failFolderPath); os.IsNotExist(err) {
		log.Err(err).Msg(errmsg.Error())
		return nil, err
	}

	ret := &FileReader{
		FolderPath:        sourceRawFolderPath,
		SuccessFolderPath: successFolderPath,
		FailFolderPath:    failFolderPath,
		allowedExtFile:    extFile,
	}

	// load files from folder path
	err := ret.loadFiles()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// loadFiles will load all files in the f.FolderPath and filter it base on the extension (if available)
func (f *FileReader) loadFiles() error {
	// read folder
	files, err := os.ReadDir(f.FolderPath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no file(s) in directory")
	}

	// sort base on file name
	// assumtion: file name is base on the date
	utils.SortFileNameAscend(files)

	// loop files to get allowed extension
	// skip if no limit of extension
	if len(f.allowedExtFile) == 0 {
		f.FileList = files
	} else {
		fileExtMap := map[string]string{}
		for _, v := range f.allowedExtFile {
			fileExtMap[v] = ""
		}

		// filter files
		for _, fl := range files {
			fileExt := strings.ReplaceAll(filepath.Ext(fl.Name()), ".", "")
			_, exist := fileExtMap[fileExt]
			if exist {
				f.FileList = append(f.FileList, fl)
			}
		}
	}

	return nil
}
