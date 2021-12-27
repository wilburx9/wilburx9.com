package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api"
)

func main() {
	api.SetUpLogrus()

	// Attempt to load config
	if err := api.LoadConfig(); err != nil {
		log.Fatalf("invalid application configuration: %s", err)
	}

	var db = api.SetUpDatabase()
	defer db.Close()

	api.ScheduleFetchAddCache(db)

	// Setup and start Http server
	s := api.SetUpServer(db)
	s.ListenAndServe()
}
