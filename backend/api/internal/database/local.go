package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"os"
	"path/filepath"
	"time"
)

const filePerm = 0644

// LocalDatabase gets and saves data to a local .json file
type LocalDatabase struct{}

// Close does nothing for the local db
func (l LocalDatabase) Close() {}

// Persist saves the data to a local .json file
func (l LocalDatabase) Write(ctx context.Context, source string, models ...Model) error {
	currentData, path, _ := l.retrieve(source)

	for _, model := range models {
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
func (l LocalDatabase) Read(ctx context.Context, source, orderBy string, limit int) ([]map[string]interface{}, UpdatedAt, error) {
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

	dir := l.getDirectory(internal.UpdatesKey)
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
