package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool    `json:"success"`
	Data    any     `json:"data,omitempty"`
	Error   *string `json:"error,omitempty"`
}

func NewResponse(w http.ResponseWriter, err *string, httpCode int, data any) {
	w.WriteHeader(httpCode)

	if httpCode == http.StatusNoContent {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response Response

	if httpCode < 300 && httpCode >= 200 && err == nil {
		response = buildPositiveResponse(data)
	} else {
		response = buildNegativeResponse(err)
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func buildPositiveResponse(data any) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func buildNegativeResponse(err *string) Response {
	return Response{
		Error:   err,
		Success: false,
	}
}
