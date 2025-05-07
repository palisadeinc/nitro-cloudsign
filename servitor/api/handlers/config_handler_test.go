package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewConfigHandler(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]string
		wantErr  bool
		wantResp []byte
	}{
		{
			name:    "valid config",
			config:  map[string]string{"key1": "value1", "key2": "value2"},
			wantErr: false,
			wantResp: func() []byte {
				r, _ := json.Marshal(map[string]string{"key1": "value1", "key2": "value2"})
				return r
			}(),
		},
		{
			name:    "empty config",
			config:  map[string]string{},
			wantErr: false,
			wantResp: func() []byte {
				r, _ := json.Marshal(map[string]string{})
				return r
			}(),
		},
		// json.Marshal for map[string]string should not error,
		// so a direct test for the error path is tricky without interface changes or mocking json.Marshal itself.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfigHandler(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConfigHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got.resp, tt.wantResp) {
				t.Errorf("newConfigHandler() got.resp = %s, want %s", got.resp, tt.wantResp)
			}
		})
	}
}

func TestConfigHandler_Handler(t *testing.T) {
	config := map[string]string{"message": "hello"}
	handler, err := newConfigHandler(config)
	if err != nil {
		t.Fatalf("failed to create config handler for tests: %v", err)
	}

	expectedBody, _ := json.Marshal(config)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedHeader map[string]string
		expectedBody   []byte
	}{
		{
			name:           "GET request - success",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedHeader: map[string]string{"Content-Type": "application/json"},
			expectedBody:   expectedBody,
		},
		{
			name:           "POST request - method not allowed",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedHeader: map[string]string{}, // Error response might not set Content-Type consistently
			expectedBody:   []byte("Method not allowed\n"),
		},
		{
			name:           "PUT request - method not allowed",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedHeader: map[string]string{},
			expectedBody:   []byte("Method not allowed\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For GET requests or others where body is nil, use http.NoBody to avoid nil pointer dereference on r.Body.Close()
			body := http.NoBody
			req, err := http.NewRequest(tt.method, "/config", body)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			httpHandler := http.HandlerFunc(handler.Handler)
			httpHandler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			for key, expectedValue := range tt.expectedHeader {
				if value := rr.Header().Get(key); value != expectedValue {
					t.Errorf("handler returned wrong header %s: got %s want %s",
						key, value, expectedValue)
				}
			}

			// For non-200 responses, http.Error adds a newline.
			// For 200 responses, we expect exact body match.
			if tt.expectedStatus == http.StatusOK {
				if !bytes.Equal(rr.Body.Bytes(), tt.expectedBody) {
					t.Errorf("handler returned unexpected body: got %s want %s",
						rr.Body.String(), string(tt.expectedBody))
				}
			} else {
				if !bytes.Equal(rr.Body.Bytes(), tt.expectedBody) {
					t.Errorf("handler returned unexpected body for error: got %q want %q",
						rr.Body.String(), string(tt.expectedBody))
				}
			}
		})
	}
}

// Test case for r.Body.Close() error path in the defer statement.
// This is a bit contrived as it's hard to make r.Body.Close() fail reliably in a test
// without a custom io.ReadCloser that always errors on Close.
// However, we can add a test to ensure the defer is executed.
type errorCloser struct{}

func (ec *errorCloser) Read(p []byte) (n int, err error) {
	return 0, nil // Doesn't matter for this test
}

func (ec *errorCloser) Close() error {
	// This is where we could log or set a flag to verify Close was called.
	// For this test, we'll just ensure the handler doesn't panic.
	// In a real scenario with logging, you might check logs.
	return errors.New("simulated close error")
}

func TestConfigHandler_Handler_BodyCloseError(t *testing.T) {
	config := map[string]string{"message": "hello"}
	handler, err := newConfigHandler(config)
	if err != nil {
		t.Fatalf("failed to create config handler for tests: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, "/config", &errorCloser{})
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.Handler)

	// We are testing if the defer r.Body.Close() is handled gracefully.
	// Log output is not easily verifiable here without more complex test setup (e.g. hooking into logrus).
	// The main thing is that the handler completes without panic and returns expected result.
	httpHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// We can add a check for the log output if we capture logs.
	// For now, we assume if it doesn't panic and returns OK, the error was logged.
}
