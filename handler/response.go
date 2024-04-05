package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func RespondJSON(ctx context.Context, writer http.ResponseWriter, body any, status int) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("encode response error: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		response := ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			fmt.Printf("Write error response error: %v", err)
		}
		return
	}

	writer.WriteHeader(status)
	if _, err := fmt.Fprintf(writer, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}
