package engine

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"stock-data-processing/fileprocessorengine/filereader"
	"stock-data-processing/model"
	"stock-data-processing/pkg/pubsub"
)

const (
	headerPub                  = "fileprocessorengine"
	topicFileProcessing string = "raw-data-ready-to-process"
)

type Engine struct {
	FileReader     *filereader.FileReader
	publisher      pubsub.Publisher
	enableMoveFile bool
}

type ProcessResult struct {
	FileCount    int
	SuccessCount int
	FailFile     []ProcessFail

	DataCount        int
	DataSuccessCount int
	FailData         []ProcessDataFail
}

type ProcessFail struct {
	FilePath string
	Err      error
}

type ProcessDataResult struct {
	DataCount    int
	SuccessCount int
	// data that skipped because of the filter
	Fail []ProcessDataFail
}

type ProcessDataFail struct {
	RawData string
	Err     error
}

func NewEngine(fileReader *filereader.FileReader, publisher pubsub.Publisher, enableMoveFile bool) (*Engine, error) {
	if fileReader == nil {
		return nil, errors.New("filereader required")
	}
	return &Engine{
		FileReader:     fileReader,
		publisher:      publisher,
		enableMoveFile: enableMoveFile,
	}, nil
}

func TestNewEngine(fileReader *filereader.FileReader, publisher pubsub.Publisher) (*Engine, error) {
	if fileReader == nil {
		return nil, errors.New("filereader required")
	}
	return &Engine{
		FileReader: fileReader,
		publisher:  publisher,
	}, nil
}

// ProcessFile will process file base on file reader
func (e *Engine) ProcessFile(ctx context.Context) ProcessResult {
	result := ProcessResult{
		FileCount: len(e.FileReader.FileList),
	}
	for _, vl := range e.FileReader.FileList {
		// read and filter file
		filePath := filepath.Join(e.FileReader.FolderPath, vl.Name())
		datas, err := e.readFile(filePath)
		if err != nil {
			log.Err(err).Msg(err.Error())
			// add fail information
			result.FailFile = append(result.FailFile, ProcessFail{
				FilePath: filePath,
				Err:      err,
			})
			if err := e.moveFileFail(vl); err != nil {
				log.Err(err).Msg(err.Error())
			}
			continue
		}
		result.SuccessCount++

		// process file
		processResult := e.processData(datas)
		// add result
		result.DataCount += processResult.DataCount
		result.DataSuccessCount += processResult.SuccessCount
		result.FailData = append(result.FailData, processResult.Fail...)
		if err := e.moveFileSuccess(vl); err != nil {
			log.Err(err).Msg(err.Error())
		}
	}

	return result
}

func (*Engine) readFile(filePath string) ([]string, error) {
	filePath = filepath.Join(filepath.Clean(filePath), "")
	fs, err := os.Open(filePath)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return []string{}, err
	}
	defer func() {
		err = fs.Close()
	}()

	result := []string{}
	scanner := bufio.NewScanner(fs)
	// loop line by line
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Err(err).Msg(err.Error())
		return []string{}, err
	}

	return result, nil
}

func (e *Engine) processData(datas []string) ProcessDataResult {
	result := ProcessDataResult{}
	result.DataCount = len(datas)
	for _, val := range datas {
		// convert string to struct
		rawData := model.Raw{}
		var err = json.Unmarshal([]byte(val), &rawData)
		if err != nil {
			log.Err(err).Msg(err.Error())
			continue
		}

		// send to kafka
		headers := pubsub.MessageHeaders{}
		headers.Add("origin", headerPub)
		id := uuid.New().String()
		err = e.publisher.Send(topicFileProcessing, id, headers, []byte(val))
		if err != nil {
			log.Err(err).Msg(err.Error())
			result.Fail = append(result.Fail, ProcessDataFail{
				RawData: val,
				Err:     err,
			})
			continue
		}
		result.SuccessCount++
	}

	return result
}

func (e *Engine) moveFileSuccess(file os.DirEntry) error {
	if !e.enableMoveFile {
		return nil
	}

	oldLocation := path.Join(e.FileReader.FolderPath, file.Name())
	newLocation := path.Join(e.FileReader.SuccessFolderPath, file.Name())
	err := os.Rename(oldLocation, newLocation)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}
	return nil
}

func (e *Engine) moveFileFail(file os.DirEntry) error {
	if !e.enableMoveFile {
		return nil
	}

	oldLocation := path.Join(e.FileReader.FolderPath, file.Name())
	newLocation := path.Join(e.FileReader.FailFolderPath, file.Name())
	err := os.Rename(oldLocation, newLocation)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}
	return nil
}
