package main

import "net/http"

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
