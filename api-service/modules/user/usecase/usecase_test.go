package usecase

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-data-processing/api-service/entities"
	mockRepo "stock-data-processing/api-service/modules/user/repository/mocks"
)

func Test_usecase_GetList(t *testing.T) {
	mockRepo := new(mockRepo.Repository)
	mockRepo.On("GetList", mock.Anything, mock.Anything).Return(
		&[]entities.User{
			{
				Name: "test1",
				Age:  13,
			},
		},
		int64(10), nil,
	)

	usecase := NewUsecase(mockRepo)
	result, err := usecase.GetList(1, 1)
	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, result)
}

func Test_usecase_Insert(t *testing.T) {
	mockRepo := new(mockRepo.Repository)
	mockRepo.On("InsertOne", mock.Anything).Return(nil)

	usecase := NewUsecase(mockRepo)
	err := usecase.Insert(entities.User{
		Name: "test-1",
		Age:  13,
	})
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err)
}

// func Test_usecase_Insert(t *testing.T) {
// 	type fields struct {
// 		repo repository.Repository
// 	}
// 	type args struct {
// 		user entities.User
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
// 			u := usecase{
// 				repo: tt.fields.repo,
// 			}
// 			if err := u.Insert(tt.args.user); (err != nil) != tt.wantErr {
// 				t.Errorf("usecase.Insert() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
