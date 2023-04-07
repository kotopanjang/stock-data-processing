package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"stock-data-processing/api-service/entities"
	"stock-data-processing/api-service/modules/user/usecase"
)

type HTTPHandler struct {
	Usecase usecase.Usecase
}

func NewHTTPHandler(router *mux.Router, usecase usecase.Usecase) {
	handler := HTTPHandler{
		Usecase: usecase,
	}

	router.HandleFunc("/get-users", handler.GetList)
	router.HandleFunc("/create-user", handler.CreateUser).Methods("POST")
}

func (u HTTPHandler) GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qq := r.URL.Query()
	_take := qq.Get("take")
	take, _ := strconv.Atoi(_take)
	_limit := qq.Get("limit")
	limit, _ := strconv.Atoi(_limit)
	result, err := u.Usecase.GetList(take, limit)
	if err != nil {
		json.NewEncoder(w).Encode(result)
	}

	json.NewEncoder(w).Encode(result)
}

func (u HTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.FormValue("name")
	_age := r.FormValue("age")
	age, _ := strconv.Atoi(_age)
	err := u.Usecase.Insert(entities.User{Name: name, Age: age})

	if err != nil {
		json.NewEncoder(w).Encode(`{"result": "Error"}`)
	}

	json.NewEncoder(w).Encode(`{"result": "ok"}`)
}
