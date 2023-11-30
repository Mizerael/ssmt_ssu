package formingquery

import (
	"log"
	"ssmt-ssu/search/synonyms"
	"strings"

	"github.com/blevesearch/bleve/v2/search/query"
)

func CreateQuery(str string) []query.Query {
	var searchQuery []query.Query
	tmpString := strings.Split(str, " ")
	for _, v := range tmpString {
		synonyms, err := synonyms.GetSynonyms(v)
		if err != nil {
			log.Fatal(err)
		}
		searchQuery = append(searchQuery, query.NewTermQuery(v))
		for _, k := range synonyms {
			searchQuery = append(searchQuery, query.NewTermQuery(k.Text))
		}
	}
	return searchQuery
}
