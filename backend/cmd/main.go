package main

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"time"
)

func main() {
	api.SetUpLogrus()

	// Attempt to load config
	if err := configs.LoadConfig("../configs"); err != nil {
		log.Fatalf("invalid application configuration: %s", err)
	}

	err := api.SetUpSentry()
	if err != nil {
		log.Fatalf("sentry.Init: failed %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))

	if err != nil {
		log.Fatalf("setting up badger failed %v", err)
		return
	}

	go api.CacheDataSources(db)

	// Setup and start Http server
	s := api.SetUpServer(db)
	s.ListenAndServe()
}
