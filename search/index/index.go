package index

import (
	"encoding/json"
	"log"
	"os"
	facewho "ssmt-ssu/search/faceWho"
	"ssmt-ssu/search/mapping"
	searchconfig "ssmt-ssu/search/searchConfig"
	"time"

	"github.com/blevesearch/bleve/v2"
)

var jsonDir = "data/"
var indexPath = "index/scp.bleve"
var batchSize = 512

type Document struct {
	Name     string `json:"name"`
	Class    string `json:"class"`
	Contains string `json:"contains"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func indexScp(i bleve.Index, conf *searchconfig.Config) error {
	dirEntries, err := os.ReadDir(jsonDir)
	checkErr(err)

	log.Printf("Indexing...")
	count := 0
	startTime := time.Now()
	batch := i.NewBatch()
	batchCount := 0

	for _, dirEntry := range dirEntries {
		filename := dirEntry.Name()
		jsonBytes, err := os.ReadFile(jsonDir + "/" + filename)
		checkErr(err)

		var jsonDoc []Document
		err = json.Unmarshal(jsonBytes, &jsonDoc)
		checkErr(err)

		for _, scp := range jsonDoc {
			docID := scp.Name
			if conf.UseSummarize {
				containsToEn, err := facewho.Translate(facewho.RuToEnApi, scp.Contains, conf)
				if err != nil {
					log.Fatal(err)
				}
				containsSummarize, err := facewho.Summarize(facewho.SummarizeApi, containsToEn, conf)
				if err != nil {
					log.Fatal(err)
				}
				scp.Contains, err = facewho.Translate(facewho.EnToRuApi, containsSummarize, conf)
				if err != nil {
					log.Fatal(err)
				}
			}
			batch.Index(docID, scp)
			batchCount++
			if batchCount >= batchSize {
				err = i.Batch(batch)
				checkErr(err)

				batch = i.NewBatch()
				batchCount = 0
				count++
			}
		}

		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	if batchCount > 0 {
		err = i.Batch(batch)
		checkErr(err)
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil
}

func OpenIndex(path string, conf *searchconfig.Config) (bleve.Index, error) {
	if conf.IndexRebuild {
		os.RemoveAll(path)
	}
	scpIndex, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index\n")
		indexMapping, err := mapping.CreateIndexMapping()
		checkErr(err)

		scpIndex, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			log.Fatal(err)
		}
		err = indexScp(scpIndex, conf)
		checkErr(err)
	} else if err != nil {

		checkErr(err)
	}
	log.Printf("Opening Index")
	return scpIndex, nil
}
