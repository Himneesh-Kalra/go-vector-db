package search

import (
	"reflect"
	"testing"

	"github.com/Himneesh-Kalra/go-vector-db/models"
)

func TestBruteSearch_TopK(t *testing.T) {
	search := BruteSearch{}

	table := map[string]models.Vector{
		"vec1": {ID: "vec1", Values: []float32{1, 0}},
		"vec2": {ID: "vec2", Values: []float32{0, 1}},
		"vec3": {ID: "vec3", Values: []float32{0.7, 0.7}},
	}

	query := []float32{1, 0}
	k := 2

	results := search.TopK(query, k, "test_table", table)

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if results[0].ID != "vec1" {
		t.Errorf("Expected first result to be vec1, got %s", results[0].ID)
	}

	if results[1].ID != "vec3" {
		t.Errorf("Expected second result to be vec3, got %s", results[1].ID)
	}
}

func TestBruteSearch_TopK_EmptyTable(t *testing.T) {
	search := BruteSearch{}
	query := []float32{1, 2, 3}
	k := 5

	results := search.TopK(query, k, "empty_table", map[string]models.Vector{})

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty table, got %d", len(results))
	}
}

func TestBruteSearch_TopK_KGreaterThanData(t *testing.T) {
	search := BruteSearch{}

	table := map[string]models.Vector{
		"a": {ID: "a", Values: []float32{1, 1}},
		"b": {ID: "b", Values: []float32{0, 1}},
	}

	query := []float32{1, 1}
	k := 5

	results := search.TopK(query, k, "small_table", table)

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	expectedIDs := []string{"a", "b"}
	actualIDs := []string{results[0].ID, results[1].ID}

	if !reflect.DeepEqual(expectedIDs, actualIDs) && !reflect.DeepEqual(reverse(expectedIDs), actualIDs) {
		t.Errorf("Expected result IDs %v, got %v", expectedIDs, actualIDs)
	}
}

func reverse(s []string) []string {
	res := make([]string, len(s))
	for i := range s {
		res[i] = s[len(s)-1-i]
	}
	return res
}
