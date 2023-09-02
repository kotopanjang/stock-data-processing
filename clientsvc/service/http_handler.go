package service

import (
	"net/http"

	"github.com/gorilla/mux"

	"stock-data-processing/clientsvc/response"
)

type HTTPHandler struct {
	Svc Service
}

func NewHTTPHandler(router *mux.Router, svc Service) {
	handler := HTTPHandler{
		Svc: svc,
	}
	router.HandleFunc("/get-summary/{stockCode}", handler.GetSummaryByStock).Methods("GET")
}

func (h HTTPHandler) GetSummaryByStock(w http.ResponseWriter, r *http.Request) {
	var resp response.Response

	pathVariable := mux.Vars(r)
	stockCode := pathVariable["stockCode"]

	resp = h.Svc.GetSummaryByStockCode(stockCode)
	response.JSON(w, resp)
}
