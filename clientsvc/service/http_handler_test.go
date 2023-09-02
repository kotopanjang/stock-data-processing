package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"stock-data-processing/clientsvc/response"
	"stock-data-processing/clientsvc/service"
	"stock-data-processing/clientsvc/service/mocks"
	"stock-data-processing/clientsvc/storage/redis"
	"stock-data-processing/model"
)

const (
	Anything = "mock.Anything"
)

func Test_GetSummaryByStock_Success(t *testing.T) {
	rds := redis.NewTestRedis()
	svc := service.NewService(rds)
	svcCollMock := mocks.NewServiceColl(t)

	hh := service.HTTPHandler{
		Svc: *svc,
	}

	successResponse := response.NewSuccessResponse(model.StockSummary{}, response.StatOK, "")
	svcCollMock.On("GetSummaryByStockCode", Anything).Return(successResponse)

	r := httptest.NewRequest(http.MethodPost, "/just/for/testing", nil)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(hh.GetSummaryByStock)

	handler.ServeHTTP(recorder, r)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected %v. got: %v", http.StatusOK, recorder.Code)
	}

	svcCollMock.AssertExpectations(t)
}

// func TestHandler_SetupRouter(t *testing.T) {
// 	rds := redis.NewTestRedis()
// 	ctrl := gomock.NewController(t)
// 	t.Cleanup(ctrl.Finish)

// 	svc := service.NewService(rds)

// 	muxRouter := mux.NewRouter()

// 	// http handler
// 	service.NewHTTPHandler(muxRouter, *svc)

// 	// test
// 	req, err := http.NewRequest("GET", "/espay/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	resp := httptest.NewRecorder()
// 	muxRouter.ServeHTTP(resp, req)
// }
