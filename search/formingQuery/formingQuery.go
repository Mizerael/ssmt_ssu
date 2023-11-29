package formingquery

import (
	"strings"

	"github.com/blevesearch/bleve/v2/search/query"
)

func CreateQuery(str string) []query.Query {
	var searchQuery []query.Query
	tmpString := strings.Split(str, " ")
	for _, v := range tmpString {
		searchQuery = append(searchQuery, query.NewTermQuery(v))
	}
	return searchQuery
}
