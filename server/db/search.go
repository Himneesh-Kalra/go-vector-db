package db

import (
	// "sort"

	"github.com/Himneesh-Kalra/go-vector-db/config"
	"github.com/Himneesh-Kalra/go-vector-db/search"
	"github.com/Himneesh-Kalra/go-vector-db/storage"
	// "github.com/Himneesh-Kalra/go-vector-db/utils"
)

var (
	Brute                        = search.BruteSearch{}
	LSH                          = search.NewLSHSearch(config.AppConfig.LSHL, config.AppConfig.LSHK, ApiQuery.Size)
	Searcher search.VectorSearch = Brute //default
)

func UseLSH(dim int) {
	l := config.AppConfig.LSHL
	k := config.AppConfig.LSHK

	LSH = search.NewLSHSearch(l, k, dim)
	Searcher = LSH
}

func UseBrute() {
	Searcher = Brute
}

func SearchTopK(query []float32, k int, table string) []search.Result {
	tableData, ok := storage.Store[table]
	ApiQuery.Size=len(query)
	if !ok {
		return nil
	}
	if lsh, ok := Searcher.(*search.LSHSearch); ok && lsh.Dim == 0 {
		lsh.Dim = len(query)
	}
	return Searcher.TopK(query, k, table, tableData)
}

type Query struct {
	Size int
}

var ApiQuery Query