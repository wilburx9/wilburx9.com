package main

import (
	"context"
	"github.com/wilburt/wilburx9.dev/backend"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"log"
)

func main() {
	// Enable date, time and line number for log outputs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Attempt to load config
	if err := common.LoadConfig("./"); err != nil {
		log.Fatalf("invalid application configuration: %s", err)
	}

	dbCtx := context.Background()
	fsClient := backend.CreateFireStoreClient(dbCtx)
	defer fsClient.Close()

	// Setup custom logger
	err := backend.SetLogger()
	if err != nil {
		log.Fatalf("sentry.Init: failed %s", err)
	}
	defer backend.CleanUpLogger()

	// Setup and start Http server
	s := backend.SetUpServer(dbCtx, fsClient)
	s.ListenAndServe()
}
