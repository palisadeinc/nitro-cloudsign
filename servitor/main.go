package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/palisadeinc/nitro-cloudsign/servitor/api"
	"github.com/palisadeinc/nitro-cloudsign/servitor/api/handlers"
	log "github.com/sirupsen/logrus"
)

const apiPort = uint32(3000)
const shutdownTimeout = 5 * time.Second

// main is the entrypoint for the servitor binary
func main() {
	ctx := context.Background()

	config, err := buildConfig()
	if err != nil {
		log.WithError(err).Fatal("Error building configuration")
	}

	if err := run(ctx, config, net.Listen); err != nil {
		log.WithError(err).Fatal("Application run failed")
	}
	log.Info("Application terminated successfully")
}

// listenFunc is a function type for providing a net.Listener, typically net.Listen.
type listenFunc func(network, address string) (net.Listener, error)

func run(ctx context.Context, config map[string]string, listen listenFunc) error {
	handler, err := handlers.Handler(config)
	if err != nil {
		return fmt.Errorf("could not create handler: %w", err)
	}

	address := fmt.Sprintf(":%d", apiPort)
	listener, err := listen("tcp", address)
	if err != nil {
		return fmt.Errorf("error listening on tcp %s: %w", address, err)
	}

	server := api.NewServer(handler)

	// Start the server. api.Server.Start is non-blocking and starts http.Serve in goroutines.
	server.Start(listener)
	log.WithContext(ctx).WithField("address", listener.Addr().String()).Info("Server started and listening")

	// Wait for context cancellation from the OS signal (e.g. CTRL+C) or other termination reason.
	<-ctx.Done()
	log.WithContext(ctx).Info("Shutdown signal received, initiating graceful shutdown...")

	// Create a context for the graceful shutdown with a timeout.
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	// Perform graceful shutdown.
	if err := server.Shutdown(shutdownCtx); err != nil {
		// This error means graceful shutdown failed (e.g., timeout exceeded).
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.WithContext(ctx).Info("Server shutdown completed successfully")
	return nil
}

func buildConfig() (map[string]string, error) {
	config := make(map[string]string)

	pairingKey, present := os.LookupEnv("PAIRING_KEY")
	if !present {
		return nil, fmt.Errorf("PAIRING_KEY is required")
	}
	config["PAIRING_KEY"] = pairingKey

	dbDataSource, present := os.LookupEnv("DB_DATA_SOURCE")
	if !present {
		return nil, fmt.Errorf("DB_DATA_SOURCE is required")
	}
	config["DB_DATA_SOURCE"] = dbDataSource

	tsmDbDataSource, present := os.LookupEnv("TSM_DB_DATA_SOURCE")
	if !present {
		return nil, fmt.Errorf("TSM_DB_DATA_SOURCE is required")
	}
	config["TSM_DB_DATA_SOURCE"] = tsmDbDataSource
	config["LOG_LEVEL"], _ = os.LookupEnv("LOG_LEVEL") // LOG_LEVEL is optional
	config["DB_DRIVER"] = "postgres"

	return config, nil
}
