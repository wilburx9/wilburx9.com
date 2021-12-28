package internal

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"google.golang.org/api/iterator"
	"os"
	"path/filepath"
	"time"
)

const (
	filePerm = 0644
)

// DbModel wraps the id method
type DbModel interface {
	Id() string
}

// UpdatedAt is a container that holds the last time a collection was updated
type UpdatedAt struct {
	T time.Time `firestore:"updated_at,serverTimestamp"`
}

// Database is the interface that wraps basic functions for saving anf retrieving data from concrete databases
type Database interface {
	Persist(source string, data ...DbModel) error
	Retrieve(source, orderBy string, limit int) ([]map[string]interface{}, UpdatedAt, error)
	Close()
}

// FirebaseFirestore gets and saves data to Firebase Firestore
type FirebaseFirestore struct {
	Client *firestore.Client
	Ctx    context.Context
}

// LocalDatabase gets and saves data to a local .json file
type LocalDatabase struct{}

// Close does nothing for the local db
func (l LocalDatabase) Close() {}

// Close closes the resources help by the Db
func (f FirebaseFirestore) Close() {
	f.Client.Close()
}

// Persist saves the data to Firebase Firestore
func (f FirebaseFirestore) Persist(source string, data ...DbModel) error {
	if len(data) == 0 {
		return fmt.Errorf("data is empty")
	}

	batch := f.Client.Batch()
	for _, m := range data {
		docId := f.Client.Collection(source).Doc(m.Id())
		batch.Set(docId, m)
	}

	batch.Set(f.Client.Collection(UpdatesKey).Doc(source), UpdatedAt{})
	_, err := batch.Commit(f.Ctx)
	if err != nil {
		return err
	}

	return nil
}

// Retrieve gets the data saved to Firebase Firestore
func (f FirebaseFirestore) Retrieve(source, orderBy string, limit int) ([]map[string]interface{}, UpdatedAt, error) {
	var data []map[string]interface{}
	q := f.Client.Collection(source).Query
	if orderBy != "" {
		q = q.OrderBy(orderBy, firestore.Desc)
	}
	if limit != 0 {
		q = q.Limit(limit)
	}

	iter := q.Documents(f.Ctx)
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
	snapshot, err := f.Client.Collection(UpdatesKey).Doc(source).Get(f.Ctx)
	if err == nil {
		snapshot.DataTo(&updatedAt)
	}
	return data, updatedAt, nil
}

// Persist saves the data to a local .json file
func (l LocalDatabase) Persist(source string, data ...DbModel) error {
	currentData, path, _ := l.retrieve(source)

	for _, model := range data {
		bytes, err := json.Marshal(model)
		if err != nil {
			continue
		}
		var m map[string]interface{}
		json.Unmarshal(bytes, &m)
		currentData[model.Id()] = m
	}

	newData, err := json.Marshal(currentData)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, newData, filePerm)
	if err != nil {
		return err
	}

	return l.persistUpdateTime(source)
}

// Retrieve gets the data saved to a local .json file. The results are  not ordered
func (l LocalDatabase) Retrieve(source, _ string, limit int) ([]map[string]interface{}, UpdatedAt, error) {
	result, _, err := l.retrieve(source)
	if err != nil {
		return nil, UpdatedAt{}, err
	}

	var data []map[string]interface{}
	for _, m := range result {
		data = append(data, m)
	}

	if limit < len(data) {
		data = data[:limit]
	}

	updatedAt, err := l.retrieveUpdatedTime(source)
	return data, updatedAt, err
}

func (l LocalDatabase) retrieve(source string) (map[string]map[string]interface{}, string, error) {
	dir := l.getDirectory("")
	path := l.getPath(dir, source)
	var result = make(map[string]map[string]interface{})

	data, err := os.ReadFile(path)
	if err != nil {
		return result, path, err
	}

	err = json.Unmarshal(data, &result)

	if err != nil {
		return result, path, err
	}

	return result, path, nil

}

func (l LocalDatabase) persistUpdateTime(source string) error {
	updatedAt, err := json.Marshal(UpdatedAt{T: time.Now()})
	if err != nil {
		return err
	}

	dir := l.getDirectory(UpdatesKey)
	path := filepath.Join(dir, fmt.Sprintf("%v", source))
	return os.WriteFile(path, updatedAt, filePerm)
}

func (l LocalDatabase) retrieveUpdatedTime(source string) (UpdatedAt, error) {
	dir := l.getDirectory("")
	path := l.getPath(dir, source)
	data, err := os.ReadFile(path)
	if err != nil {
		return UpdatedAt{}, err
	}
	var updatedAt UpdatedAt
	json.Unmarshal(data, &updatedAt)
	return updatedAt, nil
}

func (l LocalDatabase) getDirectory(subDir string) string {
	dir := filepath.Join(configs.Config.AppHome, ".db", subDir)
	os.MkdirAll(dir, os.ModePerm)
	return dir
}

func (l LocalDatabase) getPath(dir string, source string) string {
	return filepath.Join(dir, fmt.Sprintf("%v.json", source))
}
