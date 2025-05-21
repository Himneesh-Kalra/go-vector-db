package search

import (
	"testing"

	"github.com/Himneesh-Kalra/go-vector-db/models"
)

func createSampleVectors() map[string]models.Vector {
	return map[string]models.Vector{
		"v1": {ID: "v1", Values: []float32{1, 0, 0, 0}},
		"v2": {ID: "v2", Values: []float32{0, 1, 0, 0}},
		"v3": {ID: "v3", Values: []float32{0, 0, 1, 0}},
		"v4": {ID: "v4", Values: []float32{0, 0, 0, 1}},
	}
}

func TestNewLSHSearch(t *testing.T) {
	lsh := NewLSHSearch(5, 4, 4)

	if lsh.L != 5 || lsh.K != 4 || lsh.Dim != 4 {
		t.Errorf("Expected L=5, K=4, Dim=4 but got L=%d, K=%d, Dim=%d", lsh.L, lsh.K, lsh.Dim)
	}

	if len(lsh.planes) != 5 || len(lsh.planes[0]) != 4 {
		t.Errorf("Plane dimensions incorrect")
	}
}

func TestIndexTable(t *testing.T) {
	lsh := NewLSHSearch(3, 4, 4)
	vectors := createSampleVectors()
	lsh.indexTable("test", vectors)

	if _, ok := lsh.hashtables["test"]; !ok {
		t.Error("Hash table not created for table")
	}

	if _, ok := lsh.fullVectors["test"]; !ok {
		t.Error("Full vectors not stored for table")
	}
}

func TestTopKReturnsExpectedResults(t *testing.T) {
	lsh := NewLSHSearch(5, 6, 4)
	vectors := createSampleVectors()
	lsh.indexTable("test", vectors)

	query := []float32{1, 0, 0, 0} // Closest to v1
	results := lsh.TopK(query, 2, "test", vectors)

	if len(results) == 0 {
		t.Error("No results returned")
	}

	if results[0].ID != "v1" {
		t.Errorf("Expected v1 to be top result, got %s", results[0].ID)
	}
}

func TestTopKFallbackBruteForce(t *testing.T) {
	lsh := NewLSHSearch(1, 1, 4) // Min planes: likely to miss

	vectors := createSampleVectors()
	query := []float32{1, 0, 0, 0}

	// Not indexing to test fallback
	results := lsh.TopK(query, 1, "test_brute", vectors)

	if len(results) == 0 {
		t.Fatal("No results returned by brute force fallback")
	}
	if results[0].ID != "v1" {
		t.Errorf("Expected v1 from brute search, got %s", results[0].ID)
	}
}
