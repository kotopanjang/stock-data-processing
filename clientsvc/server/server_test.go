package server_test

import (
	"net/http"
	"testing"
	"time"

	"stock-data-processing/clientsvc/server"
)

func TestServer(t *testing.T) {
	httpHandler := http.NewServeMux()

	srv := server.NewServer(httpHandler, "9091")
	srv.Start()
	time.Sleep(time.Second * 1)
	srv.Close()
}
