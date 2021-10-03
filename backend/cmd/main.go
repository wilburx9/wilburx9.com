package main

import (
	"cloud.google.com/go/firestore"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api"
	"github.com/wilburt/wilburx9.dev/backend/configs"
)

func main() {
	api.SetUpLogrus()

	// Attempt to load config
	if err := api.LoadConfig(); err != nil {
		log.Fatalf("invalid application configuration: %s", err)
	}

	ctx := context.Background()
	db := createFireStoreClient(ctx)
	defer db.Close()

	api.ScheduleFetchAddCache(db, ctx)

	// Setup and start Http server
	s := api.SetUpServer(db)
	s.ListenAndServe()
}

func createFireStoreClient(ctx context.Context) *firestore.Client {
	projectId := configs.Config.GcpProjectId
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create Firestore cleint: %v", err)
	}
	return client
}
