package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"stock-data-processing/clientsvc/exception"
	"stock-data-processing/clientsvc/response"
)

func TestErrorResponse(t *testing.T) {
	resp := response.NewErrorResponse(
		exception.ErrNotFound, http.StatusNotFound, response.StatNotFound, "Resource not found",
	)

	if resp == nil {
		t.Fatal("resp should be null")
	}
	if resp.Error() != exception.ErrNotFound {
		t.Fatalf("expected %v. got: %v", exception.ErrNotFound, resp.Error())
	}
	if resp.HTTPStatusCode() != http.StatusNotFound {
		t.Fatalf("expected %v. got: %v", resp.HTTPStatusCode(), http.StatusNotFound)
	}
	if resp.Data() != nil {
		t.Fatal("resp should be null")
	}
	if resp.Status() != response.StatNotFound {
		t.Fatalf("expected %v. got: %v", resp.Status(), response.StatNotFound)
	}
}

func TestSuccessResponse(t *testing.T) {
	t.Run("when status is common ok", func(t *testing.T) {
		resp := response.NewSuccessResponse(
			nil, response.StatOK, "OK",
		)

		if resp.Error() != nil {
			t.Fatal("resp should be null")
		}
		if resp.HTTPStatusCode() != http.StatusOK {
			t.Fatalf("expected %v. got: %v", resp.HTTPStatusCode(), http.StatusOK)
		}
	})

	t.Run("when status is created", func(t *testing.T) {
		resp := response.NewSuccessResponse(
			nil, response.StatCreated, "Created",
		)

		if resp.Error() != nil {
			t.Fatal("resp should be null")
		}
		if resp.HTTPStatusCode() != http.StatusCreated {
			t.Fatalf("expected %v. got: %v", resp.HTTPStatusCode(), http.StatusCreated)
		}
	})
}

func TestRESTResponse(t *testing.T) {
	t.Run("responding json as success", func(t *testing.T) {
		recoreder := httptest.NewRecorder()
		resp := response.NewSuccessResponse(
			nil, response.StatOK, "OK",
		)
		response.JSON(recoreder, resp)

		if resp.HTTPStatusCode() != http.StatusOK {
			t.Fatalf("expected %v. got: %v", resp.HTTPStatusCode(), http.StatusOK)
		}
	})

	t.Run("responding json as error", func(t *testing.T) {
		recoreder := httptest.NewRecorder()
		resp := response.NewErrorResponse(
			exception.ErrNotFound, http.StatusNotFound, response.StatNotFound, "Resource not found",
		)
		response.JSON(recoreder, resp)

		if resp.HTTPStatusCode() != http.StatusNotFound {
			t.Fatalf("expected %v. got: %v", resp.HTTPStatusCode(), http.StatusNotFound)
		}
	})
}
