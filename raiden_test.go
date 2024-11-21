package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShoot(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		showHeaders bool
		statusCode  int
	}{
		{
			name:        "Test with headers",
			url:         "/test",
			showHeaders: true,
			statusCode:  http.StatusOK,
		},
		{
			name:        "Test without headers",
			url:         "/test",
			showHeaders: false,
			statusCode:  http.StatusOK,
		},
		{
			name:        "Test 404",
			url:         "/notfound",
			showHeaders: false,
			statusCode:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.NewServeMux()
			handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})
			handler.HandleFunc("/notfound", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})

			server := httptest.NewServer(handler)
			defer server.Close()

			Shoot(server.URL+tt.url, tt.showHeaders)
		})
	}
}
