package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	formingquery "ssmt-ssu/search/formingQuery"
	"ssmt-ssu/search/index"
	searchconfig "ssmt-ssu/search/searchConfig"
	"strings"

	"github.com/blevesearch/bleve/v2"
)

const LinkToSite = "https://scpfoundation.net"

var indexPath = "index/scp.bleve"

func main() {
	conf := searchconfig.Execute()
	reader := bufio.NewReader(os.Stdin)
	scpIndex, err := index.OpenIndex(indexPath, conf)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	var searchPhrase string
	for {
		fmt.Printf("Search:")
		searchPhrase, _ = reader.ReadString('\n')

		searchPhrase = strings.TrimRight(searchPhrase, "\n")

		query := bleve.NewDisjunctionQuery(formingquery.CreateQuery(searchPhrase, conf)...)

		searchRequest := bleve.NewSearchRequest(query)
		searchRequest.Size = conf.Count
		searchRequest.Fields = []string{"name", "contains", "class"}

		searchResult, err := scpIndex.Search(searchRequest)
		if err != nil {
			log.Fatal(err, searchResult)
		}
		fmt.Printf("searchResult:\n")
		for i, v := range searchResult.Hits {
			fmt.Printf("%d. "+LinkToSite+"/%v With tf-idf similarity: %f\n", i+1, v.Fields["name"], v.Score)
		}
	}
}
