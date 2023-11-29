package main

import (
	"fmt"
	"ssmt-ssu/search/index"

	"github.com/blevesearch/bleve/v2"
)

var indexPath = "index/scp.bleve"

func main() {
	scpIndex, err := index.OpenIndex(indexPath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// searchPhrase := "убийца скрылся во тьме"
	query := bleve.NewDisjunctionQuery(bleve.NewTermQuery("убийца"),
		bleve.NewTermQuery("скрылся"), bleve.NewTermQuery("во"), bleve.NewTermQuery("тьме"))
	query.SetBoost(2.0) // устанавливаем бустинг в 2.0

	// создаем запрос поиска с лимитом 10 результатов
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 10

	// выполняем запрос и получаем результаты
	searchResult, err := scpIndex.Search(searchRequest)
	if err != nil {
		panic(err)
	}

	fmt.Printf("searchResult: %v\n", searchResult)
}
