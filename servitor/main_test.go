package main

import (
	"context"
	"errors"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestBuildConfig(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func(t *testing.T)
		want        map[string]string
		wantErr     bool
		expectedErr string
	}{
		{
			name: "success - all required env vars present",
			setupEnv: func(t *testing.T) {
				t.Setenv("PAIRING_KEY", "test_pairing_key")
				t.Setenv("DB_DATA_SOURCE", "test_db_source")
				t.Setenv("TSM_DB_DATA_SOURCE", "test_tsm_db_source")
				t.Setenv("LOG_LEVEL", "debug")
			},
			want: map[string]string{
				"PAIRING_KEY":        "test_pairing_key",
				"DB_DATA_SOURCE":     "test_db_source",
				"TSM_DB_DATA_SOURCE": "test_tsm_db_source",
				"LOG_LEVEL":          "debug",
				"DB_DRIVER":          "postgres",
			},
			wantErr: false,
		},
		{
			name: "success - LOG_LEVEL optional and not present",
			setupEnv: func(t *testing.T) {
				t.Setenv("PAIRING_KEY", "test_pairing_key")
				t.Setenv("DB_DATA_SOURCE", "test_db_source")
				t.Setenv("TSM_DB_DATA_SOURCE", "test_tsm_db_source")
			},
			want: map[string]string{
				"PAIRING_KEY":        "test_pairing_key",
				"DB_DATA_SOURCE":     "test_db_source",
				"TSM_DB_DATA_SOURCE": "test_tsm_db_source",
				"LOG_LEVEL":          "",
				"DB_DRIVER":          "postgres",
			},
			wantErr: false,
		},
		{
			name: "error - PAIRING_KEY missing",
			setupEnv: func(t *testing.T) {
				t.Setenv("DB_DATA_SOURCE", "test_db_source")
				t.Setenv("TSM_DB_DATA_SOURCE", "test_tsm_db_source")
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "PAIRING_KEY is required",
		},
		{
			name: "error - DB_DATA_SOURCE missing",
			setupEnv: func(t *testing.T) {
				t.Setenv("PAIRING_KEY", "test_pairing_key")
				t.Setenv("TSM_DB_DATA_SOURCE", "test_tsm_db_source")
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "DB_DATA_SOURCE is required",
		},
		{
			name: "error - TSM_DB_DATA_SOURCE missing",
			setupEnv: func(t *testing.T) {
				t.Setenv("PAIRING_KEY", "test_pairing_key")
				t.Setenv("DB_DATA_SOURCE", "test_db_source")
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "TSM_DB_DATA_SOURCE is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv != nil {
				tt.setupEnv(t)
			}

			got, err := buildConfig()

			if (err != nil) != tt.wantErr {
				t.Errorf("buildConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil || err.Error() != tt.expectedErr {
					t.Errorf("buildConfig() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun_SuccessfulStartupAndShutdown(t *testing.T) {
	t.Setenv("PAIRING_KEY", "test_pk_run_success")
	t.Setenv("DB_DATA_SOURCE", "test_db_run_success")
	t.Setenv("TSM_DB_DATA_SOURCE", "test_tsm_run_success")

	config, err := buildConfig()
	if err != nil {
		t.Fatalf("buildConfig failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	errChan := make(chan error, 1)

	go func() {
		errChan <- run(ctx, config, net.Listen)
	}()

	time.Sleep(200 * time.Millisecond)

	cancel()

	select {
	case runErr := <-errChan:
		if runErr != nil {
			t.Errorf("run() returned an unexpected error: %v", runErr)
		}
	case <-time.After(shutdownTimeout + 2*time.Second):
		t.Fatal("run() did not terminate after context cancellation and shutdown timeout")
	}
}

func TestRun_ListenError(t *testing.T) {
	t.Setenv("PAIRING_KEY", "test_pk_listen_err")
	t.Setenv("DB_DATA_SOURCE", "test_db_listen_err")
	t.Setenv("TSM_DB_DATA_SOURCE", "test_tsm_listen_err")

	config, err := buildConfig()
	if err != nil {
		t.Fatalf("buildConfig failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expectedListenErr := errors.New("mock listen error")
	mockListen := func(network, address string) (net.Listener, error) {
		return nil, expectedListenErr
	}

	runErr := run(ctx, config, mockListen)

	if runErr == nil {
		t.Fatal("run() did not return an error on listen failure")
	}
	if !strings.Contains(runErr.Error(), expectedListenErr.Error()) {
		t.Errorf("run() error = %q, want error containing %q", runErr.Error(), expectedListenErr.Error())
	}
}
