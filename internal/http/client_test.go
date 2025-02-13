package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Post(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		body        []byte
		statusCode  int
		response    string
		expectError bool
	}{
		{
			name:        "valid request",
			body:        []byte("test payload"),
			statusCode:  http.StatusOK,
			response:    "response body",
			expectError: false,
		},
		{
			name:        "server error",
			body:        []byte("test payload"),
			statusCode:  http.StatusInternalServerError,
			response:    "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, "\\") {
					http.Error(w, "invalid URL", http.StatusBadRequest)
					return
				}
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.response))
			}))
			defer server.Close()

			client := &Client{ctx: context.Background(), client: server.Client()}

			resp, err := client.Post(server.URL, tc.body)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, string(resp))
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		statusCode  int
		response    string
		expectError bool
	}{
		{
			name:        "valid request",
			statusCode:  http.StatusOK,
			response:    "response body",
			expectError: false,
		},
		{
			name:        "server error",
			statusCode:  http.StatusInternalServerError,
			response:    "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.response))
			}))
			defer server.Close()

			client := &Client{ctx: context.Background(), client: server.Client()}

			resp, err := client.Get(server.URL)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, string(resp))
			}
		})
	}
}

func TestClient_do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		statusCode  int
		response    string
		expectError bool
	}{
		{
			name:        "status OK",
			statusCode:  http.StatusOK,
			response:    "response body",
			expectError: false,
		},
		{
			name:        "non-200 status",
			statusCode:  http.StatusNotFound,
			response:    "",
			expectError: true,
		},
		{
			name:        "error reading body",
			statusCode:  http.StatusOK,
			response:    "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if tc.response == "" && tc.statusCode == http.StatusOK {
					http.Error(w, "error reading body", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.response))
			}))
			defer server.Close()

			client := &Client{ctx: context.Background(), client: server.Client()}

			req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
			resp, err := client.do(req)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, string(resp))
			}
		})
	}
}
