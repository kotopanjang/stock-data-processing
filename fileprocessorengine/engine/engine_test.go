package engine

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"

	"stock-data-processing/fileprocessorengine/filereader"
	"stock-data-processing/pkg/pubsub"
	pubsub_mock "stock-data-processing/pkg/pubsub/mocks"
	// pubsub_mock "stock-data-processing/pkg/pubsub/mocks"
)

const Anything = "mock.Anything"

func Test_TestNewEngine(t *testing.T) {
	mockPubblisher := pubsub_mock.NewPublisher(t)
	fr, _ := filereader.NewFileReader("../subsetdata/raw", "../subsetdata/done/success", "../subsetdata/done/fail", []string{""})

	type args struct {
		fileReader *filereader.FileReader
		publisher  pubsub.Publisher
	}
	tests := []struct {
		name    string
		want    *Engine
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				fileReader: fr,
				publisher:  mockPubblisher,
			},
			want: &Engine{
				publisher:  mockPubblisher,
				FileReader: fr,
			},
			wantErr: false,
		},
		{
			name: "Err",
			args: args{
				fileReader: nil,
				publisher:  mockPubblisher,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TestNewEngine(tt.args.fileReader, tt.args.publisher)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEngine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NewEngine(t *testing.T) {
	mockPubblisher := pubsub_mock.NewPublisher(t)
	fr, _ := filereader.NewFileReader("../subsetdata/raw", "../subsetdata/done/success", "../subsetdata/done/fail", []string{""})

	type args struct {
		fileReader *filereader.FileReader
		publisher  pubsub.Publisher
	}
	tests := []struct {
		name    string
		want    *Engine
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				fileReader: fr,
				publisher:  mockPubblisher,
			},
			want: &Engine{
				publisher:  mockPubblisher,
				FileReader: fr,
			},
			wantErr: false,
		},
		{
			name: "Err",
			args: args{
				fileReader: nil,
				publisher:  mockPubblisher,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEngine(tt.args.fileReader, tt.args.publisher, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEngine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_ProcessFile_Success(t *testing.T) {
	mockPubblisher := pubsub_mock.NewPublisher(t)
	mockPubblisher.On("Send", Anything, Anything, Anything, Anything, Anything).Return(nil)

	fr, _ := filereader.NewFileReader("../subsetdata/test", "../subsetdata/done/success", "../subsetdata/done/fail", []string{})

	type fields struct {
		FileReader *filereader.FileReader
		publisher  pubsub.Publisher
	}
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   ProcessResult
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				FileReader: fr,
				publisher:  mockPubblisher,
			},
			want: ProcessResult{
				FileCount:        1,
				SuccessCount:     1,
				DataCount:        1,
				DataSuccessCount: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				FileReader: tt.fields.FileReader,
				publisher:  tt.fields.publisher,
			}
			if got := e.ProcessFile(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Engine.ProcessFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_readFile(t *testing.T) {
	mockPubblisher := pubsub_mock.NewPublisher(t)
	// mockPubblisher.On("Send", Anything, Anything, Anything, Anything, Anything).Return(nil)

	fr, _ := filereader.NewFileReader("../subsetdata/test", "../subsetdata/done/success", "../subsetdata/done/fail", []string{})

	type fields struct {
		FileReader *filereader.FileReader
		publisher  pubsub.Publisher
	}
	// type args struct {
	// 	filePath string
	// }
	tests := []struct {
		name   string
		fields fields
		// args    args
		want    []string
		wantErr bool
	}{
		{
			name: "readFile_success",
			fields: fields{
				publisher:  mockPubblisher,
				FileReader: fr,
			},
			want:    []string{`{"type":"P","executed_quantity":"5","order_book":"35","execution_price":"4530","stock_code":"UNVR"}`},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				FileReader: tt.fields.FileReader,
				publisher:  tt.fields.publisher,
			}

			for _, vl := range e.FileReader.FileList {
				fp := filepath.Join(e.FileReader.FolderPath, vl.Name())
				got, err := e.readFile(fp)
				if (err != nil) != tt.wantErr {
					t.Errorf("Engine.readFile() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Engine.readFile() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestEngine_processData(t *testing.T) {
	mockPubblisher := pubsub_mock.NewPublisher(t)
	mockPubblisher.On("Send", Anything, Anything, Anything, Anything, Anything).Return(nil)
	fr, _ := filereader.NewFileReader("../subsetdata/test", "../subsetdata/done/success", "../subsetdata/done/fail", []string{})

	type fields struct {
		FileReader *filereader.FileReader
		publisher  pubsub.Publisher
	}
	type args struct {
		ctx   context.Context
		datas []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ProcessDataResult
	}{
		{
			name: "processData_success",
			fields: fields{
				publisher:  mockPubblisher,
				FileReader: fr,
			},
			args: args{
				ctx:   context.Background(),
				datas: []string{`{"type":"P","executed_quantity":"5","order_book":"35","execution_price":"4530","stock_code":"UNVR"}`},
			},
			want: ProcessDataResult{
				DataCount:    1,
				SuccessCount: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				FileReader: tt.fields.FileReader,
				publisher:  tt.fields.publisher,
			}
			if got := e.processData(tt.args.datas); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Engine.processData() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestEngine_moveFileSuccess(t *testing.T) {
// 	mockPubblisher := pubsub_mock.NewPublisher(t)
// 	mockPubblisher.On("Send", Anything, Anything, Anything, Anything, Anything).Return(nil)
// 	fr, _ := filereader.NewFileReader("../subsetdata/test", "../subsetdata/done/success", "../subsetdata/done/fail", []string{})

// 	type fields struct {
// 		FileReader *filereader.FileReader
// 		publisher  pubsub.Publisher
// 	}
// 	type args struct {
// 		file os.DirEntry
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "moveFileSuccess_ok",
// 			fields: fields{
// 				publisher:  mockPubblisher,
// 				FileReader: fr,
// 			},
// 			args: args{
// 				file: fs.FileInfoToDirEntry(),
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := &Engine{
// 				FileReader: tt.fields.FileReader,
// 				publisher:  tt.fields.publisher,
// 			}
// 			if err := e.moveFileSuccess(tt.args.file); (err != nil) != tt.wantErr {
// 				t.Errorf("Engine.moveFileSuccess() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestEngine_moveFileFail(t *testing.T) {
// 	type fields struct {
// 		FileReader *filereader.FileReader
// 		publisher  pubsub.Publisher
// 	}
// 	type args struct {
// 		file os.DirEntry
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := &Engine{
// 				FileReader: tt.fields.FileReader,
// 				publisher:  tt.fields.publisher,
// 			}
// 			if err := e.moveFileFail(tt.args.file); (err != nil) != tt.wantErr {
// 				t.Errorf("Engine.moveFileFail() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
