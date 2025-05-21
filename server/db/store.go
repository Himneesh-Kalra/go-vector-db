package db

import (
	"fmt"

	"github.com/Himneesh-Kalra/go-vector-db/models"
)



var VectorStore = make(map[string]models.Vector)

func InsertVector(vec models.Vector) {
	VectorStore[vec.ID] = vec
}

func DeleteVector(id string) bool {
	if _, exists := VectorStore[id]; exists {
		delete(VectorStore, id)
		fmt.Println("Vector deleted")
		return true
	}
	fmt.Println("Vector not found")
	return false
}

func GetAllVectors() map[string]models.Vector {
	return VectorStore
}

func GetVectorValues() map[string][]float32 {
	vectors := make(map[string][]float32)
	for id, vec := range VectorStore {
		vectors[id] = vec.Values
	}
	return vectors
}
