package formingquery

import (
	"fmt"
	searchconfig "ssmt-ssu/search/searchConfig"
	"ssmt-ssu/search/spellshecker"
	"ssmt-ssu/search/synonyms"
	"strings"

	"github.com/blevesearch/bleve/v2/search/query"
)

const spellcheckPath = "https://speller.yandex.net/services/spellservice.json/checkTexts?text="

func CreateQuery(str string, conf *searchconfig.Config) []query.Query {
	var searchQuery []query.Query
	tmpString := strings.Split(str, " ")
	for _, word := range tmpString {
		potentiallyCorr, _ := spellshecker.Spellcheck(word)

		var synonims []synonyms.Syn
		maybeCorrect := false
		for _, potCorr := range potentiallyCorr {
			if word == potCorr {
				maybeCorrect = true
			}
			syn, err := synonyms.GetSynonyms(potCorr)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			} else {
				synonims = append(synonims, syn...)
			}
		}
		if !maybeCorrect {
			syn, err := synonyms.GetSynonyms(word)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			synonims = append(synonims, syn...)
		}
		wordToQuery := query.NewTermQuery(word)
		searchQuery = append(searchQuery, wordToQuery)
		fmt.Printf("synonims: %v\n", synonims)
		for _, k := range synonims {
			synonymsToQuery := query.NewTermQuery(k.Text)
			synonymsToQuery.SetBoost(conf.Bosting)
			searchQuery = append(searchQuery, query.NewTermQuery(k.Text))
		}
	}
	return searchQuery
}
