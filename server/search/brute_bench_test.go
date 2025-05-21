package search

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Himneesh-Kalra/go-vector-db/models"
)

func generateVectors(n, dim int) map[string]models.Vector {
	vectors := make(map[string]models.Vector, n)
	for i := 0; i < n; i++ {
		values := make([]float32, dim)
		for j := range values {
			values[j] = rand.Float32()
		}
		vectors[fmt.Sprintf("vec%d", i)] = models.Vector{
			ID:     fmt.Sprintf("vec%d", i),
			Values: values,
		}
	}
	return vectors
}

func BenchmarkBruteSearch_TopK_1K(b *testing.B) {
	runBenchmarkBruteSearch(b, 1000, 128, 10)
}

func BenchmarkBruteSearch_TopK_10K(b *testing.B) {
	runBenchmarkBruteSearch(b, 10_000, 128, 10)
}

func BenchmarkBruteSearch_TopK_100K(b *testing.B) {
	runBenchmarkBruteSearch(b, 100_000, 128, 10)
}

func runBenchmarkBruteSearch(b *testing.B, numVectors, dim, k int) {
	search := BruteSearch{}
	table := generateVectors(numVectors, dim)

	query := make([]float32, dim)
	for i := range query {
		query[i] = rand.Float32()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = search.TopK(query, k, "benchmark_table", table)
	}
}
