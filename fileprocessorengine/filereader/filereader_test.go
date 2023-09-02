package filereader

import (
	"os"
	"reflect"
	"testing"
)

const (
	raw     = "../subsetdata/test"
	success = "../subsetdata/done/success"
	fail    = "../subsetdata/done/fail"
)

func TestNewFileReader(t *testing.T) {
	type args struct {
		sourceRawFolderPath string
		successFolderPath   string
		failFolderPath      string
		extFile             []string
	}
	type want struct {
		fileCount int
	}
	tests := []struct {
		name    string
		args    args
		want    *want
		wantErr bool
	}{
		{
			name: "NewFileReader_ok",
			args: args{
				sourceRawFolderPath: raw,
				successFolderPath:   success,
				failFolderPath:      fail,
				extFile:             []string{},
			},
			want: &want{
				fileCount: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFileReader(tt.args.sourceRawFolderPath, tt.args.successFolderPath, tt.args.failFolderPath, tt.args.extFile)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewFileReader() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(len(got.FileList), tt.want.fileCount) {
				t.Errorf("NewFileReader() = %v, want %v", len(got.FileList), tt.want.fileCount)
			}
		})
	}
}

func TestFileReader_loadFiles(t *testing.T) {
	fr, _ := NewFileReader(raw, success, fail, []string{})
	type fields struct {
		FolderPath        string
		SuccessFolderPath string
		FailFolderPath    string
		FileList          []os.DirEntry
		allowedExtFile    []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "loadFiles_success",
			fields: fields{
				allowedExtFile:    []string{},
				FolderPath:        raw,
				SuccessFolderPath: success,
				FailFolderPath:    fail,
				FileList:          fr.FileList,
			},
		},
		{
			name: "loadFiles_success_with_extension",
			fields: fields{
				allowedExtFile:    []string{"ndjson"},
				FolderPath:        raw,
				SuccessFolderPath: success,
				FailFolderPath:    fail,
				FileList:          fr.FileList,
			},
		},
		{
			name: "loadFiles_err_empty_files",
			fields: fields{
				allowedExtFile:    []string{},
				FolderPath:        fail,
				SuccessFolderPath: success,
				FailFolderPath:    fail,
				FileList:          fr.FileList,
			},
			wantErr: true,
		},
		{
			name: "loadFiles_err_1",
			fields: fields{
				allowedExtFile:    []string{},
				FolderPath:        "",
				SuccessFolderPath: success,
				FailFolderPath:    fail,
				FileList:          fr.FileList,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileReader{
				FolderPath:        tt.fields.FolderPath,
				SuccessFolderPath: tt.fields.SuccessFolderPath,
				FailFolderPath:    tt.fields.FailFolderPath,
				FileList:          tt.fields.FileList,
				allowedExtFile:    tt.fields.allowedExtFile,
			}
			if err := f.loadFiles(); (err != nil) != tt.wantErr {
				t.Errorf("FileReader.loadFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
