package search

import "github.com/Himneesh-Kalra/go-vector-db/models"

type Result struct {
	ID     string
	Score  float32
	Values []float32
}

// Interface
type VectorSearch interface {
	TopK(query []float32, k int, tableName string, table map[string]models.Vector) []Result
}
