package main

import (
	"context"
	"net"
	"os"

	"github.com/palisadeinc/nitro-cloudsign/servitor/api"
	"github.com/palisadeinc/nitro-cloudsign/servitor/api/handlers"
	log "github.com/sirupsen/logrus"
)

const apiPort = uint32(3000)

// main is the entrypoint for the servitor binary
func main() {
	ctx := context.Background()

	config := buildConfig()

	handler, err := handlers.Handler(config)
	if err != nil {
		log.WithError(err).Error("could not create handler")
		return
	}

	tcpListener, err := net.Listen("tcp", ":3000") // TODO potentially limit to ip interface
	if err != nil {
		log.WithField("port", apiPort).WithError(err).Error("error listening on tcp")
		return
	}

	server := api.NewServer(handler)
	server.Start(tcpListener)
	log.WithContext(ctx).WithField("port", apiPort).Info("server started")

	<-ctx.Done()
}

func buildConfig() map[string]string {
	config := make(map[string]string)

	pairingKey, present := os.LookupEnv("PAIRING_KEY")
	if !present {
		log.Fatal("PAIRING_KEY is required")
	}
	config["PAIRING_KEY"] = pairingKey

	dbDataSource, present := os.LookupEnv("DB_DATA_SOURCE")
	if !present {
		log.Error("DB_DATA_SOURCE is required")
	}
	config["DB_DATA_SOURCE"] = dbDataSource

	tsmDbDataSource, present := os.LookupEnv("TSM_DB_DATA_SOURCE")
	if !present {
		log.Error("TSM_DB_DATA_SOURCE is required")
	}
	config["TSM_DB_DATA_SOURCE"] = tsmDbDataSource
	config["LOG_LEVEL"], _ = os.LookupEnv("LOG_LEVEL")
	config["DB_DRIVER"] = "postgres"

	return config
}
