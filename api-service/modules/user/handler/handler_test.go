package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-data-processing/api-service/entities"
	"stock-data-processing/api-service/modules/user/handler"
	mockUsecase "stock-data-processing/api-service/modules/user/usecase/mocks"
	"stock-data-processing/api-service/utils"
)

func Test_New_Handler(t *testing.T) {
	t.Run("should construct the customer payment method http handler", func(t *testing.T) {
		mockUsecase := new(mockUsecase.Usecase)
		router := mux.NewRouter()

		handler.NewHTTPHandler(router, mockUsecase)
	})
}

func Test_usecase_GetList(t *testing.T) {
	mockUsecase := new(mockUsecase.Usecase)
	mockUsecase.On("GetList", mock.Anything, mock.Anything).Return(
		&utils.Pagination{
			Data: []entities.User{
				{
					Name: "test1",
					Age:  13,
				},
			},
			Count:      1,
			TotalCount: 1,
			Error:      nil,
		},
		nil,
	)

	hh := handler.HTTPHandler{
		Usecase: mockUsecase,
	}

	r := httptest.NewRequest(http.MethodGet, "/just/for/testing?take=1&skip=1", nil)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(hh.GetList)
	handler.ServeHTTP(recorder, r)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockUsecase.AssertExpectations(t)
}
