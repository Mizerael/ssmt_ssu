package mapping

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/analysis/lang/ru"
	"github.com/blevesearch/bleve/v2/mapping"
)

func createMapper(str string) *mapping.FieldMapping {
	mapper := bleve.NewTextFieldMapping()
	mapper.Analyzer = str
	return mapper
}

func CreateIndexMapping() (mapping.IndexMapping, error) {
	ruTextFieldmapping := createMapper(ru.AnalyzerName)
	keywordFieldMapping := createMapper(en.AnalyzerName)

	scpMapping := bleve.NewDocumentMapping()
	scpMapping.AddFieldMappingsAt("name", keywordFieldMapping)
	scpMapping.AddFieldMappingsAt("class", ruTextFieldmapping)
	scpMapping.AddFieldMappingsAt("contains", ruTextFieldmapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("json", scpMapping)

	return indexMapping, nil
}
