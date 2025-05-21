package storage

import (
	"encoding/json"
	"fmt"

	// "io/ioutil"
	// "path"

	// "io/ioutil"
	"os"
	"path/filepath"

	"github.com/Himneesh-Kalra/go-vector-db/models"
)

var DataDir string

var Store = make(map[string]map[string]models.Vector)

func LoadAll() error {

	if _, err := os.Stat(DataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(DataDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create data directory :%v", DataDir)
		}
		fmt.Println("Created missing data directory: ", DataDir)
	}

	files, err := filepath.Glob(filepath.Join(DataDir, "*.json"))
	if err != nil {
		return err
	}
	for _, file := range files {
		table := filepath.Base(file[:len(file)-5])

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Failed to read file", file)
			continue
		}

		var vectors map[string]models.Vector

		if err := json.Unmarshal(data, &vectors); err != nil {
			fmt.Println("Failed to parse ", file)
			continue
		}
		Store[table] = vectors
	}
	return nil
}

func SaveTable(table string) error {
	if _, ok := Store[table]; !ok {
		return fmt.Errorf("no such table : %v", table)
	}
	data, err := json.MarshalIndent(Store[table], "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(DataDir, table+".json")
	fmt.Println(path)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return LoadAll()
}
