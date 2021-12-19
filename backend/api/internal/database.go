package internal

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"os"
	"path/filepath"
)

// Database is the interface that wraps basic functions for saving anf retrieving data from concrete databases
type Database interface {
	Persist(parentKey string, key string, data interface{})
	Retrieve(parentKey string, key string, result interface{}) error
	Close()
}

// FirebaseFirestore gets and saves data to Firebase Firestore
type FirebaseFirestore struct {
	Client *firestore.Client
	Ctx    context.Context
}

// Close closes the resources help by the Db
func (f FirebaseFirestore) Close() {
	f.Client.Close()
}

// Persist saves the data to Firebase Firestore
func (f FirebaseFirestore) Persist(parentKey string, key string, data interface{}) {
	mapData := map[string]interface{}{"data": data}
	_, err := f.Client.Collection(parentKey).Doc(key).Set(f.Ctx, mapData)
	if err != nil {
		log.Errorf("error: %v :: failed to write %q to %s.%s", err, mapData, parentKey, key)
	}
}

// Retrieve gets the data saved to Firebase Firestore
func (f FirebaseFirestore) Retrieve(parentKey string, key string, result interface{}) error {
	snapshot, err := f.Client.Collection(parentKey).Doc(key).Get(f.Ctx)
	if err != nil {
		log.Errorf("Failed to read from %s.%s:: %s", parentKey, key, err)
		return err
	}

	data, err := snapshot.DataAt("data")

	if err != nil {
		log.Errorf("Failed to read snapshot data at %s.%s:: %s", parentKey, key, err)
		return err
	}
	return decodeMap(data, result)
}

// LocalDatabase gets and saves data to a local .json file
type LocalDatabase struct{}

// Close does nothing for the local db
func (l LocalDatabase) Close() {}

// Persist saves the data to a local .json file
func (l LocalDatabase) Persist(parentKey string, key string, data interface{}) {
	directory := getDirectory(parentKey)

	b, err := json.Marshal(data)
	if err != nil {
		log.Errorf("error: %v :: failed to mashal data", err)
		return
	}

	os.MkdirAll(directory, os.ModePerm)
	path := getFullPath(directory, key)
	err = os.WriteFile(path, b, 0644)
	if err != nil {
		log.Errorf("error: %v :: failed to save file at %q", err, path)
	}
}

// Retrieve gets the data saved to a local .json file
func (l LocalDatabase) Retrieve(parentKey string, key string, result interface{}) error {
	path := getFullPath(getDirectory(parentKey), key)

	data, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("error: %v :: failed to read file at %q", err, path)
		return err
	}

	var dest []interface{}
	err = json.Unmarshal(data, &dest)
	if err != nil {
		log.Errorf("error: %v :: failed to unmashal data", err)
		return err
	}

	return decodeMap(dest, result)
}

func decodeMap(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(ToTimeHookFunc()),
		Result: &result,
	})
	if err != nil {
		return err
	}

	if err = decoder.Decode(input); err != nil {
		return err
	}
	return nil
}


func getDirectory(parentKey string) string {
	return filepath.Join(configs.Config.AppHome, ".db", parentKey)
}

func getFullPath(parent, name string) string {
	return filepath.Join(parent, fmt.Sprintf("%v.json", name))
}
