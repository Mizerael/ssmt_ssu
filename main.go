package main

import (
	"fmt"
	"ssmt-ssu/search/index"

	"github.com/blevesearch/bleve/v2"
)

var indexPath = "index/scp.bleve"

func main() {
	indexMapping := bleve.NewIndexMapping()
	scpIndex, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	err = index.IndexScp(scpIndex)

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
