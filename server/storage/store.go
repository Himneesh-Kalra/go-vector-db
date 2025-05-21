package storage

import (
	"fmt"
	"log"

	"github.com/Himneesh-Kalra/go-vector-db/models"
)

func InsertVector(table string, vec models.Vector) {
	if _, ok := Store[table]; !ok {
		Store[table] = make(map[string]models.Vector)
		fmt.Println("new table created")
	}
	Store[table][vec.ID] = vec
	fmt.Println("vector saved to table")
	if err := SaveTable(table); err != nil {
		log.Fatal("couldnt save table to disk", err)
	}
	fmt.Println("data saved to file")
}

func DeleteVector(table string, id string) bool {
	if tableData, ok := Store[table]; ok {
		if _, exists := tableData[id]; exists {
			delete(tableData, id)
			if err := SaveTable(table); err != nil {
				fmt.Println("error saving table after delete", err)
			}
			return true
		}
	}
	return false
}

func GetTable(table string) (map[string]models.Vector, bool) {
	data, ok := Store[table]
	return data, ok
}
