package search

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/Himneesh-Kalra/go-vector-db/models"
	"github.com/Himneesh-Kalra/go-vector-db/utils"
)

type LSHSearch struct {
	K           int // no of bits (hyperplanes) per hash
	L           int // no of hash tables
	Dim         int // dimension of vectors
	hashtables  map[string][]map[string]map[string]bool
	planes      [][][]float32
	fullVectors map[string]map[string]models.Vector
	rng         *rand.Rand
}

func NewLSHSearch(L int, K int, dim int) *LSHSearch {
	src := rand.NewSource(time.Now().UnixNano())

	if dim <= 0 {
		panic("Invalid dimension : must be positive")
	}

	lsh := LSHSearch{
		L:           L,
		K:           K,
		Dim:         dim,
		hashtables:  make(map[string][]map[string]map[string]bool),
		fullVectors: make(map[string]map[string]models.Vector),
		rng:         rand.New(src),
	}

	lsh.generatePlanes()
	return &lsh
}

func (l *LSHSearch) indexTable(table string, vectors map[string]models.Vector) {

	l.fullVectors[table] = vectors

	tables := make([]map[string]map[string]bool, l.L)

	for i := range tables {
		tables[i] = make(map[string]map[string]bool)
	}

	for id, vec := range vectors {
		normalizedVec := normalizeVector(vec.Values)
		for ti := 0; ti < l.L; ti++ {
			sig := make([]byte, l.K)
			for hi := 0; hi < l.K; hi++ {
				if utils.Dot(normalizedVec, l.planes[ti][hi]) >= 0 {
					sig[hi] = '1'
				} else {
					sig[hi] = '0'
				}
			}
			key := string(sig)

			if tables[ti][key] == nil {
				tables[ti][key] = make(map[string]bool)
			}
			tables[ti][key][id] = true
		}
	}
	l.hashtables[table] = tables
}

func (l *LSHSearch) TopK(query []float32, k int, tableName string, table map[string]models.Vector) []Result {

	fmt.Println("Before normalizing: ", query)

	query = normalizeVector(query)

	fmt.Println("After normalizing: ", query)

	if l.Dim != len(query) {
		l.Dim = len(query)
		fmt.Println("New dim is ", l.Dim)
		l.generatePlanes()
	}
	// if l.Dim == 0 {
	// 	l.Dim = len(query)
	// 	l.generatePlanes()
	// }

	if _, ok := l.hashtables[tableName]; !ok {
		fmt.Printf("Indexing table %s with lsh ... \n", tableName)
		l.indexTable(tableName, table)
	}

	hashTables, ok := l.hashtables[tableName]
	if !ok {
		fmt.Println("No hashtable found for table : ", tableName)
		return nil
	}

	fmt.Printf("Using lsh search for table %s . \n", tableName)

	candidates := make(map[string]bool)

	for ti := 0; ti < l.L; ti++ {
		sig := make([]byte, l.K)

		for hi := 0; hi < l.K; hi++ {
			dot := utils.Dot(query, l.planes[ti][hi])
			fmt.Printf("Dot product with plane [%d][%d]  :  %f\n", ti, hi, dot)
			if dot >= 0 {
				sig[hi] = '1'
			} else {
				sig[hi] = '0'
			}
			fmt.Printf("Signature = [ %c ] \n", sig[hi])
		}
		key := string(sig)

		if ids, ok := hashTables[ti][key]; ok {
			fmt.Printf("Table %d, Signatuure %s, fouund %d candidates \n", ti, key, len(ids))
			for id := range ids {
				candidates[id] = true
			}
		}
	}

	if len(candidates) == 0 {
		fmt.Println("No candidates found using LSH , falling back to brute search")
		bruteSearch := BruteSearch{}
		return bruteSearch.TopK(query, k, tableName, table)
	}
	results := make([]Result, 0, len(candidates))
	for id := range candidates {
		vec, ok := table[id]
		if !ok {
			continue
		}
		score := utils.CosineSimilarity(query, vec.Values)
		results = append(results, Result{
			ID:     id,
			Score:  score,
			Values: vec.Values,
		})
	}
	sort.Slice(results, func(i int, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > k {
		return results[:k]
	}
	return results

}

func (l *LSHSearch) generatePlanes() {
	l.planes = make([][][]float32, l.L)
	for i := 0; i < l.L; i++ {

		l.planes[i] = make([][]float32, l.K)
		for j := 0; j < l.K; j++ {

			vec := make([]float32, l.Dim)
			totalNorm := float64(0)

			for d := 0; d < l.Dim; d++ {
				v := l.rng.NormFloat64()
				vec[d] = float32(v)
				totalNorm += v * v
			}

			if totalNorm > 1e-9 {
				normFactor := float32(1.0 / math.Sqrt(totalNorm))
				for d := 0; d < l.Dim; d++ {
					vec[d] *= normFactor
				}
			}

			l.planes[i][j] = vec

		}
	}
}

func normalizeVector(vec []float32) []float32 {
	norm := float32(0)
	for _, v := range vec {
		norm += v * v
	}

	norm = float32(math.Sqrt(float64(norm)))

	if norm < 1e-6 {
		return vec
	}
	normalized := make([]float32, len(vec))
	for i, v := range vec {
		normalized[i] = v / norm
	}
	return normalized
}
