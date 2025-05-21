package search

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Himneesh-Kalra/go-vector-db/models"
)

// Helper to generate large number of random vectors
func generateRandomVectors(n int, dim int) map[string]models.Vector {
	vectors := make(map[string]models.Vector, n)
	for i := 0; i < n; i++ {
		vec := make([]float32, dim)
		for j := 0; j < dim; j++ {
			vec[j] = rand.Float32()
		}
		vectors[fmt.Sprintf("vec%d", i)] = models.Vector{
			ID:     fmt.Sprintf("vec%d", i),
			Values: vec,
		}
	}
	return vectors
}

// Benchmark TopK on a large table using LSH
func BenchmarkTopK_LSH(b *testing.B) {
	dim := 128
	numVectors := 100_000
	query := make([]float32, dim)
	for i := range query {
		query[i] = rand.Float32()
	}

	// Prepare LSH and index data
	lsh := NewLSHSearch(10, 12, dim)
	vectors := generateRandomVectors(numVectors, dim)
	lsh.indexTable("benchmark_table", vectors)

	b.ResetTimer() // resets the timer so setup code doesn't affect the benchmark

	for i := 0; i < b.N; i++ {
		_ = lsh.TopK(query, 10, "benchmark_table", vectors)
	}
}

func BenchmarkTopK_BruteForce(b *testing.B) {
	dim := 128
	numVectors := 100_000
	query := make([]float32, dim)
	for i := range query {
		query[i] = rand.Float32()
	}

	// Set very low L and K to force LSH to miss and fallback
	lsh := NewLSHSearch(1, 1, dim)
	vectors := generateRandomVectors(numVectors, dim)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = lsh.TopK(query, 10, "brute_table", vectors)
	}
}
