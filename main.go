package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	facewho "ssmt-ssu/search/faceWho"
	formingquery "ssmt-ssu/search/formingQuery"
	"ssmt-ssu/search/index"
	searchconfig "ssmt-ssu/search/searchConfig"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/joho/godotenv"
)

const LinkToSite = "https://scpfoundation.net"

var indexPath = "index/scp.bleve"

func main() {
	conf := searchconfig.Execute()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	conf.YandexApiKey = os.Getenv("YandexApiKey")
	conf.HuggingfaceApiKey = os.Getenv("HuggingfaceApiKey")

	reader := bufio.NewReader(os.Stdin)
	scpIndex, err := index.OpenIndex(indexPath, conf)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	var searchPhrase string
	for {
		fmt.Printf("Search:")
		searchPhrase, _ = reader.ReadString('\n')

		searchPhrase = strings.TrimRight(searchPhrase, "\n")

		query := bleve.NewDisjunctionQuery(formingquery.CreateQuery(searchPhrase, conf)...)

		searchRequest := bleve.NewSearchRequest(query)
		searchRequest.Size = conf.Count * 5
		searchRequest.Fields = []string{"name", "contains", "class"}

		searchResult, err := scpIndex.Search(searchRequest)
		if err != nil {
			log.Fatal(err, searchResult)
		}
		result := facewho.Simmilarity(facewho.SimilarirtyApi, searchPhrase, searchResult, conf)
		fmt.Printf("searchResult:\n")
		for i, v := range result {
			fmt.Printf("%d. "+LinkToSite+"/%v With similarity: %f\n", i+1, v.Key, v.Value)
		}
	}
}
