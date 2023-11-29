package main

import (
	"fmt"
	formingquery "ssmt-ssu/search/formingQuery"
	"ssmt-ssu/search/index"
	searchconfig "ssmt-ssu/search/searchConfig"

	"github.com/blevesearch/bleve/v2"
)

var indexPath = "index/scp.bleve"

func main() {
	conf := searchconfig.Execute()

	scpIndex, err := index.OpenIndex(indexPath, conf)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	var searchPhrase string

	fmt.Printf("Search: ")
	fmt.Scanf("%s\n", &searchPhrase)

	query := bleve.NewDisjunctionQuery(formingquery.CreateQuery(searchPhrase)...)
	query.SetBoost(float64(conf.Bosting))

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = conf.Count

	searchResult, err := scpIndex.Search(searchRequest)
	if err != nil {
		panic(err)
	}

	fmt.Printf("searchResult: %v\n", searchResult)
}
