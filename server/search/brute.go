package search

import (
	"sort"

	"github.com/Himneesh-Kalra/go-vector-db/models"
	// "github.com/Himneesh-Kalra/go-vector-db/storage"
	"github.com/Himneesh-Kalra/go-vector-db/utils"
)

type BruteSearch struct{}

func (b BruteSearch) TopK(query []float32, k int, tableName string, table map[string]models.Vector) []Result {
	results := make([]Result, 0)

	for id, vec := range table {
		score := utils.CosineSimilarity(query, vec.Values)
		results = append(results, Result{ID: id, Score: score, Values: vec.Values})
	}
	sort.Slice(results, func(i int, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) < k {
		return results
	}
	return results[:k]
}
