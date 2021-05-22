package main

import (
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend"
	"github.com/wilburt/wilburx9.dev/backend/common"
)

func main() {
	backend.SetUpLogrus()

	// Attempt to load config
	if err := common.LoadConfig("./"); err != nil {
		log.Fatalf("invalid application configuration: %s", err)
	}

	err := backend.SetUpSentry()
	if err != nil {
		log.Fatalf("sentry.Init: failed %s", err)
	}
	defer backend.CleanUpLogger()

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))

	if err != nil {
		log.Fatalf("setting up badger failed %v", err)
		return
	}

	go backend.CacheDataSources(db)

	// Setup and start Http server
	s := backend.SetUpServer(db)
	s.ListenAndServe()
}
