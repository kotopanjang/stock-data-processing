package service

import (
	"net/http"
	"stock-data-processing/clientsvc/response"
	"stock-data-processing/clientsvc/storage/redis"
	"stock-data-processing/model"
)

type ServiceColl interface {
	GetSummaryByStockCode(stockCode string) (resp response.Response)
}

type Service struct {
	Store *redis.RedisStorage
}

func NewService(redis *redis.RedisStorage) *Service {
	return &Service{
		Store: redis,
	}
}

func (s Service) GetSummaryByStockCode(stockCode string) (resp response.Response) {
	sum := &model.StockSummary{}
	err := s.Store.GetByKey(stockCode, sum)
	if err != nil {
		return response.NewErrorResponse(err, http.StatusInternalServerError, response.StatUnexpectedError, "error get data from redis")
	}

	if sum.StockCode == "" {
		return response.NewErrorResponse(err, http.StatusNotFound, response.StatNotFound, "no data found")
	}

	return response.NewSuccessResponse(sum, response.StatOK, "")
}
