// COSINE DISTANCE ALGORITHM
package utils

import (
	"fmt"
	"math"
)

func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		fmt.Printf("Dimension mismatch. wrong query size. len(a): %d \t len(b): %d \n \n", len(a), len(b))
		return 0
	}
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
