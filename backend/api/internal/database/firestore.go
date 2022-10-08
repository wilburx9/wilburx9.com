package database

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"google.golang.org/api/iterator"
)

// FirebaseFirestore gets and saves data to Firebase Firestore
type FirebaseFirestore struct {
	Client *firestore.Client
}

// Persist saves the data to Firebase Firestore
func (f FirebaseFirestore) Write(ctx context.Context, source string, models ...Model) error {
	if len(models) == 0 {
		return fmt.Errorf("models is empty")
	}

	batch := f.Client.Batch()
	for _, m := range models {
		docId := f.Client.Collection(source).Doc(m.Id())
		batch.Set(docId, m)
	}

	batch.Set(f.Client.Collection(internal.UpdatesKey).Doc(source), UpdatedAt{})
	_, err := batch.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Retrieve gets the data saved to Firebase Firestore
func (f FirebaseFirestore) Read(ctx context.Context, source, orderBy string, limit int) ([]map[string]interface{}, UpdatedAt, error) {
	var data []map[string]interface{}
	q := f.Client.Collection(source).Query
	if orderBy != "" {
		q = q.OrderBy(orderBy, firestore.Desc)
	}
	if limit != 0 {
		q = q.Limit(limit)
	}

	iter := q.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, UpdatedAt{}, err
		}
		data = append(data, doc.Data())
	}

	var updatedAt UpdatedAt
	snapshot, err := f.Client.Collection(internal.UpdatesKey).Doc(source).Get(ctx)
	if err == nil {
		snapshot.DataTo(&updatedAt)
	}
	return data, updatedAt, nil
}

// Close closes the resources help by the Db client
func (f FirebaseFirestore) Close() {
	f.Client.Close()
}
