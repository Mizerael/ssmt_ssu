package mapping

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/ru"
	"github.com/blevesearch/bleve/v2/mapping"
)

func CreateIndexMapping() (mapping.IndexMapping, error) {
	ruTextFieldmapping := bleve.NewTextFieldMapping()

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = ru.AnalyzerName

	scpMapping := bleve.NewDocumentMapping()
	scpMapping.AddFieldMappingsAt("class", keywordFieldMapping)
	scpMapping.AddFieldMappingsAt("contains", ruTextFieldmapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("scp", scpMapping)

	return indexMapping, nil
}
